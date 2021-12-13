mkdir -p $1/1 $1/2

printf "package main\n\n" >> $1/1/main.go
printf "import (\n" >> $1/1/main.go
printf "\t\"io/ioutil\"\n" >> $1/1/main.go
printf "\t\"os\"\n" >> $1/1/main.go
printf ")\n\n" >> $1/1/main.go
printf "func main() {\n" >> $1/1/main.go
printf "\tinput, err := ioutil.ReadAll(os.Stdin)\n" >> $1/1/main.go
printf "\tif err != nil {\n" >> $1/1/main.go
printf "\t\tpanic(err)\n" >> $1/1/main.go
printf "\t}\n" >> $1/1/main.go
printf "\tinput = input[:len(input)-1] //Remove last \\n\n" >> $1/1/main.go
printf "}" >> $1/1/main.go

vim $1/test.txt
vim $1/input.txt
vim $1/1/main.go
cp $1/1/main.go $1/2/
vim $1/2/main.go
