import express, { json, response } from "express";

const app = express();
app.use(express.json());
app.get("/method3", async (req, res) => {
    const arr = [];
    for (let i = 0; i < 10; i++) {
        const request = await fetch("http://3.105.180.131:3000/method2_response", {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            }
        });
        arr.push(request);
    }

    Promise.all(arr)
      .then(async (responses) => {
        const dataPromise = responses.map(response => response.json())
        const data = await Promise.all(dataPromise)
        res.json(data)
      })
      .catch((error) => console.error(error))
    
})


app.listen(8080, () => {
    console.log("port is opened at 8080")
})