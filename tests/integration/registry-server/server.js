const express = require('express');
const app = express();
const PORT = 3000;

app.get('/:name/:version', (req, res) => {
    const { name, version } = req.params;
    res.send(`Name: ${name}, Version: ${version}`);
});

app.listen(PORT, () => {
    console.log(`Server is running on http://localhost:${PORT}`);
});
