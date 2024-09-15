'use strict';
const log4js = require('log4js');
const logger = log4js.getLogger('BasicNetwork');
const bodyParser = require('body-parser');
const http = require('http')
const util = require('util');
const express = require('express')
const app = express();
const expressJWT = require('express-jwt');
const jwt = require('jsonwebtoken');
const bearerToken = require('express-bearer-token');
const cors = require('cors');
const constants = require('./config/constants.json')
const bcrypt = require("bcrypt");
const collection = require("./config");
//const { jwtDecode } = require("jwt-decode");

const host = process.env.HOST || constants.host;
const port = process.env.PORT || constants.port;


const helper = require('./app/helper')
const invoke = require('./app/invoke')
const qscc = require('./app/qscc')
const query = require('./app/query')

app.options('*', cors());
app.use(cors());
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({
    extended: false
}));
// set secret variable
app.set('secret', 'thisismysecret');
app.use(expressJWT({
    secret: 'thisismysecret'
}).unless({
    path: ['/users', '/users/login', '/register', '/signup', '/login', '/users/list', '/user/username']
}));
app.use(bearerToken());

logger.level = 'debug';


app.use((req, res, next) => {
    logger.debug('New req for %s', req.originalUrl);
    if (req.originalUrl.indexOf('/users') >= 0 || req.originalUrl.indexOf('/users/login') >= 0 || req.originalUrl.indexOf('/register') >= 0 || req.originalUrl.indexOf('/signup') >= 0 || req.originalUrl.indexOf('/login') >= 0 || req.originalUrl.indexOf('/users/list') >= 0 ||  req.originalUrl.indexOf('/user/username') >= 0) {
        return next();
    }
    var token = req.token;
    jwt.verify(token, app.get('secret'), (err, decoded) => {
        if (err) {
            console.log(`Error ================:${err}`)
            res.send({
                success: false,
                message: 'Failed to authenticate token. Make sure to include the ' +
                    'token returned from /users call in the authorization header ' +
                    ' as a Bearer token'
            });
            return;
        } else {
            req.username = decoded.username;
            req.orgname = decoded.orgName;
            logger.debug(util.format('Decoded from JWT token: username - %s, orgname - %s', decoded.username, decoded.orgName));
            return next();
        }
    });
});

var server = http.createServer(app).listen(port, function () { console.log(`Server started on ${port}`) });
logger.info('****************** SERVER STARTED ************************');
logger.info('***************  http://%s:%s  ******************', host, port);
server.timeout = 240000;

function getErrorMessage(field) {
    var response = {
        success: false,
        message: field + ' field is missing or Invalid in the request'
    };
    return response;
}

// Register and enroll user
app.post('/users', async function (req, res) {
    var username = req.body.username;
    var orgName = req.body.orgName;
    logger.debug('End point : /users');
    logger.debug('User name : ' + username);
    logger.debug('Org name  : ' + orgName);
    if (!username) {
        res.json(getErrorMessage('\'username\''));
        return;
    }
    if (!orgName) {
        res.json(getErrorMessage('\'orgName\''));
        return;
    }

    var token = jwt.sign({
        exp: Math.floor(Date.now() / 1000) + parseInt(constants.jwt_expiretime),
        username: username,
        orgName: orgName

        // username: req.body.username,
        // password: req.body.password,
        // fullname: req.body.fullname,
        // email: req.body.email,
        // role: req.body.role,
        // orgName: req.body.orgName,
        // centre: req.body.centre,
    }, app.get('secret'));

    let response = await helper.getRegisteredUser(username, orgName, true);

    logger.debug('-- returned from registering the username %s for organization %s', username, orgName);
    if (response && typeof response !== 'string') {
        logger.debug('Successfully registered the username %s for organization %s', username, orgName);
        response.token = token;
        //mise a jour du status du user dans mongoose   //const username = req.params.username;        // const { status } = "enrolled";
        try {
            const updatedUser = await collection.findOneAndUpdate({ username: username }, { $set: { status: "enrolled" } }, { new: true });
            console.log("status :", updatedUser);
        } catch (error) {
            console.log("Echec mise a jour status ", error.message);
        }
        res.json(response);
    } else {
        logger.debug('Failed to register the username %s for organization %s with::%s', username, orgName, response);
        res.json({ success: false, message: response });
    }

});

// Register and enroll user
app.post('/register', async function (req, res) {
    var username = req.body.username;
    var orgName = req.body.orgName;
    logger.debug('End point : /users');
    logger.debug('User name : ' + username);
    logger.debug('Org name  : ' + orgName);
    if (!username) {
        res.json(getErrorMessage('\'username\''));
        return;
    }
    if (!orgName) {
        res.json(getErrorMessage('\'orgName\''));
        return;
    }

    var token = jwt.sign({
        exp: Math.floor(Date.now() / 1000) + parseInt(constants.jwt_expiretime),
        username: username,
        orgName: orgName
    }, app.get('secret'));

    console.log(token)

    let response = await helper.registerAndGerSecret(username, orgName);


    logger.debug('-- returned from registering the username %s for organization %s', username, orgName);
    if (response && typeof response !== 'string') {

        logger.debug('Successfully registered the username %s for organization %s', username, orgName);
        response.token = token;

        // const updatedUsers = await collection.updateOne({username: username}, {$set:{status: "enrolled"}});
        // console.log("status Update :", updatedUsers);

        res.json(response);

    } else {
        logger.debug('Failed to register the username %s for organization %s with::%s', username, orgName, response);
        res.json({ success: false, message: response });
    }

});

// Login and get jwt
app.post('/users/login', async function (req, res) {
    var username = req.body.username;
    var orgName = req.body.orgName;
   // var role = req.body.role;
    logger.debug('End point : /users');
    logger.debug('User name : ' + username);
    logger.debug('Org name  : ' + orgName);
    if (!username) {
        res.json(getErrorMessage('\'username\''));
        return;
    }
    if (!orgName) {
        res.json(getErrorMessage('\'orgName\''));
        return;
    }
    // if (!role) {
    //     res.json(getErrorMessage('\'role\''));
    //     return;
    // }

    let isUserRegistered = await helper.isUserRegistered(username, orgName);

    if (isUserRegistered) {
        try {
            const check = await collection.findOne({ username: req.body.username });
            console.log(check);
    
            if (!check) {
                res.json({ success: false, message: `User with username ${username} is not registered with ${orgName}, Please register first.` });
            }
            //compare the hash password from database with the plain text if user is autorized
           
                var token = jwt.sign({
                    exp: Math.floor(Date.now() / 1000) + parseInt(constants.jwt_expiretime),
                    username: username,
                    orgName: orgName,
                    fullname: check.fullname,
                    centre: check.centre,
                    role: check.role,
                }, app.get('secret'));
            
                console.log(token);
    
                res.json({ success: true, message: { token: token } });
        } catch {
                res.json({ success: false, message: `User with username ${username} is not registered with ${orgName}, Please register first.` });
        }


    

    } else {
        res.json({ success: false, message: `User with username ${username} is not registered with ${orgName}, Please register first.` });
    }
///////////////////////////////////////////////////////////////////////
  

   
});


// Invoke transaction on chaincode on target peers
app.post('/channels/:channelName/chaincodes/:chaincodeName', async function (req, res) {
    try {
        logger.debug('==================== INVOKE ON CHAINCODE ==================');
        var peers = req.body.peers;
        var chaincodeName = req.params.chaincodeName;
        var channelName = req.params.channelName;
        var fcn = req.body.fcn;
        var args = req.body.args;
        var transient = req.body.transient;
        console.log(`Transient data is ;${transient}`)
        logger.debug('channelName  : ' + channelName);
        logger.debug('chaincodeName : ' + chaincodeName);
        logger.debug('fcn  : ' + fcn);
        logger.debug('args  : ' + args);
        if (!chaincodeName) {
            res.json(getErrorMessage('\'chaincodeName\''));
            return;
        }
        if (!channelName) {
            res.json(getErrorMessage('\'channelName\''));
            return;
        }
        if (!fcn) {
            res.json(getErrorMessage('\'fcn\''));
            return;
        }
        if (!args) {
            res.json(getErrorMessage('\'args\''));
            return;
        }

        let message = await invoke.invokeTransaction(channelName, chaincodeName, fcn, args, req.username, req.orgname, transient);
        console.log(`message result is : ${message}`)

        const response_payload = {
            result: message,
            error: null,
            errorData: null
        }
        res.send(response_payload);

    } catch (error) {
        const response_payload = {
            result: null,
            error: error.name,
            errorData: error.message
        }
        res.send(response_payload)
    }
});

app.get('/channels/:channelName/chaincodes/:chaincodeName', async function (req, res) {
    try {
        logger.debug('==================== QUERY BY CHAINCODE ==================');

        var channelName = req.params.channelName;
        var chaincodeName = req.params.chaincodeName;
        console.log(`chaincode name is :${chaincodeName}`)
        let args = req.query.args;
        let fcn = req.query.fcn;
        let peer = req.query.peer;

        logger.debug('channelName : ' + channelName);
        logger.debug('chaincodeName : ' + chaincodeName);
        logger.debug('fcn : ' + fcn);
        logger.debug('args : ' + args);

        if (!chaincodeName) {
            res.json(getErrorMessage('\'chaincodeName\''));
            return;
        }
        if (!channelName) {
            res.json(getErrorMessage('\'channelName\''));
            return;
        }
        if (!fcn) {
            res.json(getErrorMessage('\'fcn\''));
            return;
        }
        if (!args) {
            if( fcn != 'queryAllDocs'){
               res.json(getErrorMessage('\'args\''));
            return; 
            }            
        }
        if( fcn != 'queryAllDocs'){
        console.log('args==========', args);
        args = args.replace(/'/g, '"');
        args = JSON.parse(args);
        logger.debug(args);
        }

        let message = await query.query(channelName, chaincodeName, args, fcn, req.username, req.orgname);

        const response_payload = {
            result: message,
            error: null,
            errorData: null
        }

        res.send(response_payload);
    } catch (error) {
        const response_payload = {
            result: null,
            error: error.name,
            errorData: error.message
        }
        res.send(response_payload)
    }
});

app.get('/qscc/channels/:channelName/chaincodes/:chaincodeName', async function (req, res) {
    try {
        logger.debug('==================== QUERY BY CHAINCODE ==================');

        var channelName = req.params.channelName;
        var chaincodeName = req.params.chaincodeName;
        console.log(`chaincode name is :${chaincodeName}`)
        let args = req.query.args;
        let fcn = req.query.fcn;
        // let peer = req.query.peer;

        logger.debug('channelName : ' + channelName);
        logger.debug('chaincodeName : ' + chaincodeName);
        logger.debug('fcn : ' + fcn);
        logger.debug('args : ' + args);

        if (!chaincodeName) {
            res.json(getErrorMessage('\'chaincodeName\''));
            return;
        }
        if (!channelName) {
            res.json(getErrorMessage('\'channelName\''));
            return;
        }
        if (!fcn) {
            res.json(getErrorMessage('\'fcn\''));
            return;
        }
        if (!args) {
            res.json(getErrorMessage('\'args\''));
            return;
        }
        console.log('args==========', args);
        args = args.replace(/'/g, '"');
        args = JSON.parse(args);
        logger.debug(args);

        let response_payload = await qscc.qscc(channelName, chaincodeName, args, fcn, req.username, req.orgname);

        // const response_payload = {
        //     result: message,
        //     error: null,
        //     errorData: null
        // }

        res.send(response_payload);
    } catch (error) {
        const response_payload = {
            result: null,
            error: error.name,
            errorData: error.message
        }
        res.send(response_payload)
    }
});

//******************************************************************************* */


app.get("/login", (req, res) => {
    res.render("Login");
});

app.get("/signup", (req, res) => {
    res.render("SignUp");
});

//Register User
app.post("/signup", async (req, res) => {
    const data = {
        username: req.body.username,
        password: req.body.password,
        fullname: req.body.fullname,
        email: req.body.email,
        role: req.body.role,
        orgName: req.body.orgName,
        centre: req.body.centre,
        status: "new"
    }

    var token = jwt.sign({
        exp: Math.floor(Date.now() / 1000) + parseInt(constants.jwt_expiretime),
        username: data.username,
        orgName: data.orgName,
        fullname: data.fullname,
        centre: data.centre,
        role: data.role

    }, app.get('secret'));

    console.log(token)



    let response = await helper.registerSimple(data.username, data.orgName, token);

    logger.debug('-- returned from registering the username %s for organization %s', data.username, data.orgName);
    if (response && typeof response !== 'string') {
        //Check if the User already exist in the DB
        const existingUser = await collection.findOne({ username: data.username });
        if (existingUser) {
            res.json({ success: false, message: "User already exists. Please choose a different username." });
        } else {
            //Hash the password using bcrypt
            const saltRounds = 10; //Nbre of salt round for bcrypt
            const hashedPassword = await bcrypt.hash(data.password, saltRounds);
            data.password = hashedPassword; // Replace the password with original pwd
            const userdata = await collection.insertMany(data);
            console.log(userdata);
            logger.debug('Successfully registered the username %s for organization %s', data.username, data.orgName);
            response.token = token;
            res.json(response);
            //    res.json({ success: "true", message: "success" });
        }
    } else {
        logger.debug('Failed to register the username %s for organization %s with::%s', data.username, data.orgName, response);
        res.json({ success: false, message: response });
    }
});



//Login User
app.post("/login", async (req, res) => {
    var username = req.body.username;
    var orgName = req.body.orgName;
    logger.debug('End point : /');
    logger.debug('User name : ' + username);
    logger.debug('Org name  : ' + orgName);
    if (!username) {
        res.json(getErrorMessage('\'username\''));
        return;
    }
    if (!orgName) {
        res.json(getErrorMessage('\'orgName\''));
        return;
    }   

    try {
        const check = await collection.findOne({ username: req.body.username });
        console.log(check);

        if (!check) {
            res.json({ message: "Username cannot found. Please register first !" });
        }
        //compare the hash password from database with the plqin text
        const isPasswordMatch = await bcrypt.compare(req.body.password, check.password);
        if (isPasswordMatch) {
            // res.json({ success: "true", role: check.role });
            var token = jwt.sign({
                exp: Math.floor(Date.now() / 1000) + parseInt(constants.jwt_expiretime),
                username: username,
                orgName: orgName,
                fullname: check.fullname,
                centre: check.centre,
                role: check.role,
            }, app.get('secret'));
        
            console.log(token);

            res.json({ success: "true", message: { token: token } });
        } else {
            res.json({ message: "Wrong password!. Please try again." });
        }

    } catch {
        //res.json({ success: "false", message: "Wrong details!. Please try again." });
        res.json({ success: false, message: `User with username ${username} is not registered with ${orgName}, Please register first.` });
    }
});


// Route GET pour récupérer la liste des utilisateurs au format JSON
app.get("/users/list", async (req, res) => {
    var username = "cni"; // req.body.username;
    var orgName = "Org2"; //req.body.orgName;
    logger.debug('End point : /');
    logger.debug('User name : ' + username);
    logger.debug('Org name  : ' + orgName);
    if (!username) {
        res.json(getErrorMessage('\'username\''));
        return;
    }
    if (!orgName) {
        res.json(getErrorMessage('\'orgName\''));
        return;
    }

    var token = jwt.sign({
        exp: Math.floor(Date.now() / 1000) + parseInt(constants.jwt_expiretime),
        username: username,
        orgName: orgName
        // fullname: check.fullname,
        // centre: check.centre,
        // role: check.role
    }, app.get('secret'));


    // Récupérer la liste des utilisateurs depuis la base de données
    const users = await collection.find({}); // Récupère tous les utilisateurs, vous pouvez ajouter des conditions si nécessaire
    res.json(users);

});

// Route GET pour récupérer les infos d'un user by username
app.post("/user/username", async (req, res) => {   
         
   //const username = req.params.username;
 var username = req.body.username;
    var orgName = req.body.orgName;
   // var role = req.body.role;
    logger.debug('End point : /user/username');
    logger.debug('User name : ' + username);
    logger.debug('Org name  : ' + orgName);
    if (!username) {
        res.json(getErrorMessage('\'username\''));
        return;
    }
    if (!orgName) {
        res.json(getErrorMessage('\'orgName\''));
        return;
    }
    // if (!role) {
    //     res.json(getErrorMessage('\'role\''));
    //     return;
    // }
console.log("username ", username, "orgName ", orgName)
    var token = jwt.sign({
        exp: Math.floor(Date.now() / 1000) + parseInt(constants.jwt_expiretime),
        username: username,
        orgName: orgName,
        // fullname: check.fullname,
        // centre: check.centre,
        // role: check.role
    }, app.get('secret'));

   // console.log("token ", token)
    let response = await helper.registerSimple(username, orgName, token);
    // Récupérer les datas des utilisateurs depuis la base de données

    try {
       // const user = await collection.findOne({ username: username });
        const user = await collection.findOne({ username: username });
        console.log(user);
        if (!user) {
            return res.status(404).json({ message: "Utilisateur non trouvé" });
        }

        // Retournez les informations de profil de l'utilisateur      // response.token = token;    // res.status(200).json(user);
     
        res.json(user);
    } catch (err) {
        console.error("Erreur lors de la récupération des informations de l'utilisateur :", err);
       // res.status(500).json({ message: "Erreur lors de la récupération des informations de l'utilisateur" });
        res.json(response);
    }      

});