if [ "$#" -ne 2 ]; then
    echo "[@] Usage: ./new_eth_project <directory> <name>"
    exit
fi

mkdir "$PWD/$1/$2"
touch "$PWD/$1/$2/project.ethylene"
ID = $(($RANDOM + ($RANDOM << 16) + ($RANDOM << 32) + ($RANDOM << 48))

printf "//id: This project's UUID\nid($ID);\n\n//root: This project's folder\nroot(\"root\");" >> "$PWD/$1/$2/project.ethylene"
