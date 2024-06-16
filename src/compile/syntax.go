package compile

import (
	"geth-cody/ast"
	"geth-cody/compile/data"
	"geth-cody/compile/lexer"
	"geth-cody/compile/syntax"
	"geth-cody/io"
)

func getProject(fp io.Path, projectFiles *data.AsyncSet[*io.Project]) (*io.Project, io.Error) {
	mpp, err := io.MainProjectPath(fp)
	if err != nil {
		return nil, err
	}

	project, ok := projectFiles.GetString(mpp.String())
	if !ok {
		project, err = io.NewProject(mpp)
		if err != nil {
			return nil, err
		}
		projectFiles.Set(project)
	}

	return project, nil
}

func parseFile(fp, mainDirPath io.Path, project *io.Project, unvisitedPaths *data.Chan[io.Path], symbolMap syntax.SymbolMap) (ast.File, io.Error) {
	input, err := fp.Read()
	if err != nil {
		return nil, err
	}

	tokens, err := lexer.NewLexer(input, fp).Lex()
	if err != nil {
		return nil, err
	}

	return syntax.NewParser(tokens, project, fp, mainDirPath, unvisitedPaths, symbolMap).Parse()
}

func Syntax(mainFilePath io.Path) (syntax.SymbolMap, io.Error) {
	mainDirPath := mainFilePath.Dir()
	visitedPaths, projectFiles := data.NewAsyncSet[io.Path](), data.NewAsyncSet[*io.Project]()

	unvisitedPaths := data.NewChan[io.Path](io.Env.BufferSize)
	unvisitedPaths.Send(mainFilePath)

	fileEntries := map[string]ast.File{}
	symbolEntries, closeSymbolEntries := data.RunUntilClosed(io.Env.BufferSize,
		func(se syntax.FileEntry) {
			fileEntries[se.Path.String()] = se.File
		},
	)

	symbolMap := syntax.SymbolMap{
		Projects: projectFiles.Map(),
		Files:    fileEntries,
	}

	unvisitedPaths.AsyncForEach(io.Env.BufferSize, io.Env.ThreadCount, func(fp io.Path) io.Error {
		// skip if file is already visited
		if !visitedPaths.Set(fp) {
			return nil
		}

		project, err := getProject(fp, projectFiles)
		if err != nil {
			return err
		}

		file, err := parseFile(fp, mainDirPath, project, unvisitedPaths, symbolMap)
		if err != nil {
			return err
		}

		symbolEntries <- syntax.FileEntry{
			Path: fp,
			File: file,
		}
		if unvisitedPaths.Size() == 0 {
			unvisitedPaths.Close()
		}

		return nil
	})

	closeSymbolEntries()
	return symbolMap, nil
}
