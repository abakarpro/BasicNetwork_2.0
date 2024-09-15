const mongoose = require("mongoose");
const connect = mongoose.connect("mongodb://localhost:27017/login");


//check database connected or not
connect.then(() => {
    console.log("Database connected successfuly");
})
.catch(() => {
    console.log("Database cannot be connected, please check the mongod server");
});


// Create a schema
const LoginSchema = new mongoose.Schema({  
    username:   {
        type: String,
        required: true
    },
    fullname:   {
        type: String,
        required: true
    },
    email:   {
        type: String,
        required: true
    },
    password:   {
        type: String,
        required: true
    },
    orgName:   {
        type: String,
        required: true
    },
    centre:   {
        type: String,
        required: true
    },
    role:   {
        type: String,
        required: true
    },
    status:   {
        type: String,
        required: true
    }
});

//collection Part  Model
const collection = new mongoose.model("users", LoginSchema);

module.exports = collection;