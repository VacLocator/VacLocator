package cli

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"

	"github.com/nheingit/learnGo/blockchain"
	"github.com/nheingit/learnGo/blockchain/wallet"
)

type CommandLine struct{}

//Dictionary MaxAmmount[distrito]
//Blockchain Wallet[distrito].CurrentGente -> Sale de EDADES
//If Current >

var MaxAmmount = make(map[string]int)
var CurrentAmmount = make(map[string]int)
var WalletsAddress = make(map[string]string)
var NameToDistr = make(map[string]string)

//printUsage will display what options are availble to the user
func (cli *CommandLine) printUsage() {
	fmt.Println("Usage: ")
	fmt.Println("getbalance -address ADDRESS - get balance for ADDRESS")
	fmt.Println("createblockchain -address ADDRESS creates a blockchain and rewards the mining fee")
	fmt.Println("printchain - Prints the blocks in the chain")
	fmt.Println("send -from FROM -to TO -amount AMOUNT - Send amount of coins from one address to another")
	fmt.Println("createwallet - Creates a new wallet")
	fmt.Println("listaddresses - Lists the addresses in the wallet file")
}

//validateArgs ensures the cli was given valid input
func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		//go exit will exit the application by shutting down the goroutine
		// if you were to use os.exit you might corrupt the data
		runtime.Goexit()
	}
}

//printChain will display the entire contents of the blockchain
func (cli *CommandLine) printChain() {
	chain := blockchain.ContinueBlockChain("")
	defer chain.Database.Close()
	iterator := chain.Iterator()

	for {
		block := iterator.Next()
		fmt.Printf("Previous hash: %x\n", block.PrevHash)
		fmt.Printf("hash: %x\n", block.Hash)
		pow := blockchain.NewProofOfWork(block)
		fmt.Printf("Pow: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
		// This works because the Genesis block has no PrevHash to point to.
		if len(block.PrevHash) == 0 {
			break
		}
	}
}

//listAddresses will list all addresses in the wallet file
func (cli *CommandLine) listAddresses() {
	wallets, _ := wallet.CreateWallets()
	addresses := wallets.GetAllAddresses()

	for _, address := range addresses {
		fmt.Println(address)
	}

}

//createWallet will create a wallet in the wallet file
func (cli *CommandLine) createWallet() {
	wallets, _ := wallet.CreateWallets()
	address := wallets.AddWallet()
	wallets.SaveFile()

	fmt.Printf("New address is: %s\n", address)

}

func readCSVFromUrl(url string) ([][]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	reader := csv.NewReader(resp.Body)
	reader.Comma = ';'
	data, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (cli *CommandLine) CheckPopulation(ls []string) []string {
	outls := []string{}

	for _, s := range ls {
		if MaxAmmount[s] > CurrentAmmount[NameToDistr[s]] {
			outls = append(outls, s)
		}

	}

	return outls
}

//createWallet will create a wallet in the wallet file
func (cli *CommandLine) populate() {

	urlUsers := "https://raw.githubusercontent.com/VacLocator/VacLocator/dev/Data/personas_data_aleatoria.csv"
	dataUsers, err := readCSVFromUrl(urlUsers)
	if err != nil {
		panic(err)
	}

	header_size := 7
	dataUsers = dataUsers[header_size:]

	urlDistr := "https://raw.githubusercontent.com/VacLocator/VacLocator/dev/Data/Centros_vacuna_distritos.csv"
	dataDistr, errD := readCSVFromUrl(urlDistr)
	if errD != nil {
		panic(errD)
	}

	dataDistr = dataDistr[header_size:]

	fmt.Printf("Finished Fetching Data")

	newChain := blockchain.InitBlockChain("Vac")
	newChain.Database.Close()

	fmt.Println("Finished creating chain")

	for idx, row := range dataDistr {
		val0, _ := strconv.Atoi(row[6])
		MaxAmmount[row[2]] = val0
		NameToDistr[row[2]] = row[5]
		if idx == 10 {
			break
		}
	}

	for idx, row := range dataUsers {
		val0, _ := strconv.Atoi(row[6])
		CurrentAmmount[row[2]] += val0

		if idx == 10 {
			break
		}
	}
}

//Creates a blockchain and awards address the coinbase
func (cli *CommandLine) createBlockChain(address string) {
	newChain := blockchain.InitBlockChain(address)
	newChain.Database.Close()
	fmt.Println("Finished creating chain")
}

func (cli *CommandLine) getBalance(address string) {
	chain := blockchain.ContinueBlockChain(address)
	defer chain.Database.Close()

	balance := 0
	UTXOs := chain.FindUTXO(address)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of %s: %d\n", address, balance)
}

func (cli *CommandLine) send(from, to string, amount int) {
	chain := blockchain.ContinueBlockChain(from)
	defer chain.Database.Close()

	tx := blockchain.NewTransaction(from, to, amount, chain)

	chain.AddBlock([]*blockchain.Transaction{tx})
	fmt.Println("Success!")

}

//Run will start up the command line
func (cli *CommandLine) Run() {
	cli.validateArgs()

	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	listAddressesCmd := flag.NewFlagSet("listaddresses", flag.ExitOnError)
	populateCmd := flag.NewFlagSet("populate", flag.ExitOnError)
	checkPopulationCmd := flag.NewFlagSet("checkPopulation", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	switch os.Args[1] {
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "listaddresses":
		err := listAddressesCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createwallet":
		err := createWalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "populate":
		err := populateCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "checkpopulation":
		err := checkPopulationCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			runtime.Goexit()
		}
		cli.getBalance(*getBalanceAddress)
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			runtime.Goexit()
		}
		cli.createBlockChain(*createBlockchainAddress)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			runtime.Goexit()
		}

		cli.send(*sendFrom, *sendTo, *sendAmount)
	}
	if listAddressesCmd.Parsed() {
		cli.listAddresses()
	}
	if createWalletCmd.Parsed() {
		cli.createWallet()
	}

	if populateCmd.Parsed() {
		cli.populate()
	}
	if checkPopulationCmd.Parsed() {
		ls := []string{}
		cli.CheckPopulation(ls)
	}

}
