# GO ERC20 TRACKER

GO program to fetch top 15 INST holders using ERC20 Transfer events.

## Getting Started

### Dependencies

* Go-etherium

### Executing program

* Set env variables for ethereum endpoint.
* Execute using GO
```
go run main.go
```

##Workflow
1. I have connected ethereum network using infura endpoints.
2. Using the provided contract address, I am filtering the logs.
3. I copied abi file and parse it using given abi parser in go-ethereum for reading the logs.
4. I iterate over logs and used parsed abi to extract data and stored it in structure array.
5. I sort it and slice to give 15 INST holder.
6. Thus, printing the outcome.
