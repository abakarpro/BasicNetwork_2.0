const { Gateway, Wallets, TxEventHandler, GatewayOptions, DefaultEventHandlerStrategies, TxEventHandlerFactory } = require('fabric-network');
const fs = require('fs');
const path = require("path")
const log4js = require('log4js');
const logger = log4js.getLogger('BasicNetwork');
const util = require('util')

const helper = require('./helper')

const invokeTransaction = async (channelName, chaincodeName, fcn, args, username, org_name, transientData) => {
    try {
        logger.debug(util.format('\n============ invoke transaction on channel %s ============\n', channelName));

        // load the network configuration
        // const ccpPath =path.resolve(__dirname, '..', 'config', 'connection-org1.json');
        // const ccpJSON = fs.readFileSync(ccpPath, 'utf8')
        const ccp = await helper.getCCP(org_name) //JSON.parse(ccpJSON);

        // Create a new file system based wallet for managing identities.
        const walletPath = await helper.getWalletPath(org_name) //path.join(process.cwd(), 'wallet');
        const wallet = await Wallets.newFileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the user.
        let identity = await wallet.get(username);
        if (!identity) {
            console.log(`An identity for the user ${username} does not exist in the wallet, so registering user`);
            await helper.getRegisteredUser(username, org_name, true)
            identity = await wallet.get(username);
            console.log('Run the registerUser.js application before retrying');
            return;
        }

        

        const connectOptions = {
            wallet, identity: username, discovery: { enabled: true, asLocalhost: true },
            eventHandlerOptions: {
                commitTimeout: 100,
                strategy: DefaultEventHandlerStrategies.NETWORK_SCOPE_ALLFORTX
            }
            // transaction: {
            //     strategy: createTransactionEventhandler()
            // }
        }

        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();
        await gateway.connect(ccp, connectOptions);

        // Get the network (channel) our contract is deployed to.
        const network = await gateway.getNetwork(channelName);

        const contract = network.getContract(chaincodeName);

        let result
        let message;
        if (fcn === "createDoc" || fcn === "createPrivateDocImplicitForOrg1"
            || fcn == "createPrivateDocImplicitForOrg2") {
            result = await contract.submitTransaction(fcn, args[0], args[1], args[2], args[3], args[4], args[5], 
                args[6], args[7], args[8], args[9], args[10],  args[11], args[12],  args[13],
                args[14], args[15], args[16],  args[17], args[18], args[19],  args[20], args[21], 
                args[22],  args[23], args[24], args[25], args[26]);
            message = `Successfully added the doc asset with key ${args[0]}`

        } else if (fcn === "changeDocCentre") {
            result = await contract.submitTransaction(fcn, args[0], args[1]);
            message = `Successfully changed doc Centre with key ${args[0]}`
        } else if (fcn === "changeDocStatus") {
            result = await contract.submitTransaction(fcn, args[0], args[1], args[2], args[3]);
            message = `Successfully changed doc Status with key ${args[0]}`
        } else if (fcn === "updateDocFather") {
            result = await contract.submitTransaction(fcn, args[0], args[1], args[2], args[3], args[4], args[5], 
                args[6], args[7], args[8]);
            message = `Successfully changed doc Status with key ${args[0]}`
        } else if (fcn == "createPrivateDoc" || fcn =="updatePrivateData") {
            console.log(`Transient data is : ${transientData}`)
            let docData = JSON.parse(transientData)
            console.log(`doc data is : ${JSON.stringify(docData)}`)
            let key = Object.keys(docData)[0]
            const transientDataBuffer = {}
            transientDataBuffer[key] = Buffer.from(JSON.stringify(docData.doc))
            result = await contract.createTransaction(fcn)
                .setTransient(transientDataBuffer)
                .submit()
            message = `Successfully submitted transient data`
        }
        else {
            return `Invocation require either createDoc or changeDocCentre as function but got ${fcn}`
        }

        await gateway.disconnect();

        result = JSON.parse(result.toString());

        let response = {
            message: message,
            result
        }

        return response;


    } catch (error) {

        console.log(`Getting error: ${error}`)
        return error.message

    }
}

exports.invokeTransaction = invokeTransaction;