## FLETA

FLETA Daemon repository of FLETA project official repositories.<br/>
In the Beta Testnet the FLETA Daemon, Observer Daemon, Block Explorer, and a DApp Example are provided.<br/>
In this repository, an explanation and contents about the full node feature provided by FLETA Daemon is being written.<br/>

## Related Repositories
* You can find the Node at https://github.com/fletaio/cmd/tree/master/node.
* You can find the Formulator at https://github.com/fletaio/cmd/tree/master/formulator.
* You can find the Observer at https://github.com/fletaio/cmd/tree/master/observer.
* You can find the Block Explorer at https://github.com/fletaio/block_explorer.

## Building the source

### Prerequisites
Golang 1.11 or later<br/>
C compiler(gcc, mingw etc.)<br/>

### Step
Open the Command Line and enter the command below.

```
$ go get -u github.com/fletaio/cmd
$ go install github.com/fletaio/cmd
```

You can find fleta binary in GOPATH/bin directory

## Environment

### Command line and flags
Configuration and specification are continuously changing, and currently all configurations are provided via configuration file.<br/>
When the changing possibility of configurations and command sets are reduced, command line and flags will be provided.<br/>

### System requirements

| Resource | Recommended | Minimum |
|:---------:|---------------|---------|
|CPU|6 core for 10,000 TPS|Minimum 2 core for 3,000 TPS|
|RAM|4 GB|Minimum 1 GB|
|Disk|16 GB|Minimum 4 GB|

All the requirements mentioned on this document are written based on FLETA’s full node. All requirements can be changed at the release stage.<br/>
Currently we can touch 10,000 TPS in 6 cores with 40 servers in global regions, and now we develop our system to touch 20,000 TPS that is objective of the project.

### Configuration

In current stage, only the minimal configuration for the execution is being provided.<br/>
According to the progress of the development, there can be additions/deletions/changes at any time.<br/>
Please see details in related repositories.

## Running on FLETA Beta test network

The FLETA Daemon is enabled by default starting with FLETA beta testnet<br/>
If you want to run the FLETA binary without any modification, then you are immediately connected to the Beta test network.<br/>

## Running on your own network

In order to build your own network, the Genesis and observer nodes must first be constructed.<br/>
When building the Genesis, use the files https://github.com/fletaio/cmd/node/init.go, https://github.com/fletaio/cmd/formulator/init.go and https://github.com/fletaio/cmd/observer/init.go<br/>
GenesisContext provides the initial accounts and formulators and at least one formulator is necessary, so you should write same codes to all init.go files<br/>
The below code is included in the init.go file. 
```
func initGenesisContextData(act *data.Accounter, tran *data.Transactor) (*data.ContextData, error) {
	loader := data.NewEmptyLoader(act.ChainCoord(), act, tran)
	ctd := data.NewContextData(loader, nil)

	acg := &accCoordGenerator{}
	addSingleAccount(loader, ctd, common.MustParsePublicHash("3Zmc4bGPP7TuMYxZZdUhA9kVjukdsE2S8Xpbj4Laovv"), common.NewAddress(acg.Generate(), 0), "fleta.io")
	addFormulator(loader, ctd, common.MustParsePublicHash("gDGAcf9V9i8oWLTeayoKC8bdAooNVaFnAeQKy4CsUB"), common.MustParseAddress("3CUsUpvEK"), "fleta.io.fr00001")
	...
	return ctd, nil
}
```

After creating the private key, the two sections of code from above `3Zmc4bGPP7TuMYxZZdUhA9kVjukdsE2S8Xpbj4Laovv` and `2xASBuEWw6LcQGjYxeGZH9w1DUsEDt7fvUh8p3auxyN` must be modified.<br/>
Below are the hex private key and base58 public hash with which to modify the two sections of code above.<br/>
If you follow the instructions, then you can immediately run the test.<br/>
Change the hex private key to the FormulatorKey that is located in FLETA daemon’s Configuration file. And change the address located in the comment to the FormulatorAddress.<br/>

| Hex Private Key | Base58 Public Hash |
|:--------------:|---------------------|
|`30ea36fdc9ecb0b4c2a9eb5a82f8f5784f278409fb5cfa53cf99bbed9ce49265`|`2CQBhmtferf2qWDjqSnEE3f1ECimj4Lck2CxndgqEVq`|
|`f6d94eb4131bda99277f3bc44fc498527ecd43177872a2b58ee7008225037a18`|`4D5m6ssnsf3NxJmqKg7PpwoyG2PdMNPAuQjpB8ZKjDo`|

If you modify the following above table then it will appear as below.
```
func initGenesisContextData(loader data.Loader, ctd *data.ContextData) error {
	acg := &accCoordGenerator{}
	addSingleAccount(loader, ctd, common.MustParsePublicHash("2CQBhmtferf2qWDjqSnEE3f1ECimj4Lck2CxndgqEVq"), common.NewAddress(acg.Generate(), 0), "private") // base58 address : 1111111
	addFormulator(loader, ctd, common.MustParsePublicHash("4D5m6ssnsf3NxJmqKg7PpwoyG2PdMNPAuQjpB8ZKjDo"), common.NewAddress(acg.Generate(), 0), "private00001") // base58 address : 3CUsUpvEK
}
```

By using the modified file and the component of the init.go, You can build your own observers.<br/>
To execute your private network, you need at least 5 observer nodes and 1 formulator.<br/>
The formulator must be built in a location where it is able to connect to observer network.<br/>
Also you should creates your own key and public hashes for your observer nodes<br/>
When building according to this document, create the file named config.toml and submit with the information.<br/>

Now if you run the FLETA daemon then a message saying you have successfully connected with Observer will appear.
## License

All codes under this repository are licensed under the [GNU Lesser General Public License v3.0](https://www.gnu.org/licenses/lgpl-3.0.en.html), also included in our repository in the `LICENSE` file.

