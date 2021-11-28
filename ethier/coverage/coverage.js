const SolCover = require('solidity-coverage/lib/api');
const { spawn, exec } = require('child_process');
const fs = require('fs');
const ganache = require('ganache-cli');
const path = require('path');
const utils = require('solidity-coverage/plugins/resources/plugin.utils');

async function runCoverage(){
    const api = new SolCover();

    const skip = {
        './contracts': ['random/PRNG.sol'],
        './tests': []
    };

    const toInstrument = [];
    const toCopy = [];
    for (dir in skip) {
        const {targets, skipped} = utils.assembleFiles({"contractsDir": dir}, skip[dir]);
        toInstrument.push(...targets);
        toCopy.push(...skipped);
    }
    const instrumented = api.instrument(toInstrument);

    for (src of instrumented.concat(toCopy)) {
        const dir = path.dirname(src.relativePath);
        const base = path.basename(src.relativePath, '.sol');
        const dest = path.join(
            path.dirname(src.relativePath),
            `.${path.basename(src.relativePath, '.sol')}.cover.sol`
        );
        fs.writeFileSync(dest, src.source);
    }
    
    const env = JSON.parse(JSON.stringify(process.env)); // cloned
    env.ETHIER_COVERAGE = true; // so much for const! JavaScript is stupid
    
    const gen = spawn('npm', ['run', 'generate'], {env: env});
    gen.stdout.on('data', data => console.log(data.toString()));
    gen.stderr.on('data', data => console.error(data.toString()));

    const processClosed = (proc) => {
        return new Promise((resolve, reject) => {
            proc.on('close', (code) => {
                if (code == 0) {
                    resolve();
                } else {
                    reject();
                }
            });
        });
    }

    await processClosed(gen);

    
    const coverReport = /\[ETHIER_COVERAGE\](.+?)\[ETHIER_COVERAGE\]/;
    const test = exec('npm run testverbose', {env: env}, (err, stdout, stderr) => {
        for (let line of stdout.split("\n")) {
            const report = line.match(coverReport);
            if (!report) {
                console.info(line);
                continue;
            }
            console.info(line);
            
            for (const [site, counter] of Object.entries(JSON.parse(report[1]))) {
                if (site in api.instrumenter.instrumentationData) {
                    api.instrumenter.instrumentationData[site].hits += counter.hits;
                }
            }
        }
        console.info(api.instrumenter.instrumentationData);
    });
    
    processClosed(test)
        .then(() => console.info("done"))
        .catch(() => console.error("error"));
}

runCoverage();