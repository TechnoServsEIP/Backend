const express = require('express')
const Discord = require("discord.js")
const dotenv = require('dotenv')
const app = express()
const client = new Discord.Client()
const port = 3000


app.use(express.json())
dotenv.config()
let channel
client.on('ready', () => {
    console.log('Ready!')
    channel = client.channels.cache.get(process.env.CHANNEL)
    console.log(channel)
})

app.post("/log", (req, res)=> {
    channel.send(req.body.logFile)
    .then()
    .catch(console.error)
    channel.send(req.body.message)
    .then()
    .catch(console.error)
    res.status(200).send()
})
client.login(process.env.DISCORD_TOKEN)

app.listen(port, ()=> {
    console.log(`listen on port ${port}`)
})
