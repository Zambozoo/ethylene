package compile

import (
	"geth-cody/ast"
	"geth-cody/compile/data"
	"geth-cody/compile/lexer"
	"geth-cody/compile/syntax"
	"geth-cody/io"
	"geth-cody/io/path"
)

func getProject(fp path.Path, projectFiles *data.AsyncSet[*path.Project]) (*path.Project, io.Error) {
	mpp, err := fp.MainProjectPath()
	if err != nil {
		return nil, err
	}

	project, ok := projectFiles.GetString(mpp.String())
	if !ok {
		project, err = path.NewProject(mpp)
		if err != nil {
			return nil, err
		}
		projectFiles.Set(project)
	}

	return project, nil
}

func parseFile(fp, mainDirPath path.Path, pathProvider path.Provider, project *path.Project, unvisitedPaths *data.Chan[path.Path], symbolMap syntax.SymbolMap) (ast.File, io.Error) {
	input, err := fp.Read()
	if err != nil {
		return nil, err
	}

	tokens, err := lexer.NewLexer(input, fp).Lex()
	if err != nil {
		return nil, err
	}

	return syntax.NewParser(tokens, project, fp, mainDirPath, pathProvider, unvisitedPaths, symbolMap).Parse()
}

func Syntax(pathProvider path.Provider, mainFilePath path.Path) (syntax.SymbolMap, io.Error) {
	mainDirPath := mainFilePath.Dir()
	visitedPaths, projectFiles := data.NewAsyncSet[path.Path](), data.NewAsyncSet[*path.Project]()
	symbolMap := syntax.SymbolMap{
		Projects: projectFiles.Map(),
		Files:    map[string]ast.File{},
	}

	unvisitedPaths := data.NewChan[path.Path](io.Env.BufferSize)
	unvisitedPaths.Send(mainFilePath)
	print(mainFilePath.String())

	symbolEntries, closeSymbolEntries := data.RunUntilClosed(io.Env.BufferSize,
		func(se syntax.FileEntry) {
			symbolMap.Files[se.Path.String()] = se.File
		},
	)
	defer closeSymbolEntries()

	if err := unvisitedPaths.AsyncForEach(io.Env.BufferSize, io.Env.ThreadCount, func(fp path.Path) io.Error {
		// skip if file is already visited
		if !visitedPaths.Set(fp) {
			return nil
		}

		project, err := getProject(fp, projectFiles)
		if err != nil {
			return err
		}

		file, err := parseFile(fp, mainDirPath, pathProvider, project, unvisitedPaths, symbolMap)
		if err != nil {
			return err
		}

		symbolEntries <- syntax.FileEntry{
			Path: fp,
			File: file,
		}

		return nil
	}); err != nil {
		return syntax.SymbolMap{}, err
	}
	return symbolMap, nil
}
