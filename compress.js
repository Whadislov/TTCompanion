// This program globally compresses a .wasm into a .wasm.br
// Compresses .wasm into .wasm.br
// Adjust version in index.html
// Deletes .wasm
// type node compress.js in the terminal

const fs = require('fs').promises;
const path = require('path');
const brotli = require('brotli');

const inputFile = 'wasm/TTCompanion.wasm';
const outputFile = 'wasm/TTCompanion.wasm.br';

const compressBrotli = async () => {
    try {
        const wasmBuffer = await fs.readFile(inputFile);
        // Compress with Brotli
        const compressed = brotli.compress(wasmBuffer, {
            mode: 0, // 0 generic compression (default), 1 for text
            quality: 11, // compression level (0 to 11, 11 is the max)
            lgwin: 22 // window size (default 22)
        });

        // Verify that the compression succeed
        if (compressed === null) {
            console.error('Failure');
            process.exit(1);
        }

        await fs.writeFile(outputFile, compressed);
                console.log(`File compressed with Brotli: ${outputFile}`);
    } catch (err) {
                console.error('Error during compression:', err);
                process.exit(1);
    }
};

const deleteWasmFile = async () => {
    try {
        await fs.unlink(inputFile);
        console.log('TTCompanion.wasm has been deleted');
    } catch (err) {
        console.error('Erreur. TTCompanion.wasm has not been deleted :', err);
    }
};

const updateVersion = async (newVersion) => {
    try {
        const filePath = path.join(__dirname, 'wasm/index.html'); 
        let content = await fs.readFile(filePath, 'utf8');

        const oldLine = '"application-version">v0.0.1<';
        const newLine = `"application-version">v${newVersion}<`;
        content = content.replace(oldLine, newLine);

        await fs.writeFile(filePath, content, 'utf8');
        console.log(`Updated version : ${newVersion}`);
    } catch (err) {
        console.error('Error. Version has not been updated :', err);
    }
};

// Main func
const compress = async () => {
    // Get parameter passed in the terminal. First parameter after "node script.js"
    const newVersion = process.argv[2];
    // Check if parameter is given
    if (!newVersion) {
        console.error("Error : give a version as argument (ex. node script.js 1.2.3)");
        process.exit(1);
    }

    await compressBrotli();
    await deleteWasmFile();
    await updateVersion(newVersion);
    console.log('Scripts are done');
};

compress().catch(err => console.error('Error :', err));