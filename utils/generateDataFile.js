// Quick and dirty utility to create the server/server_data/appData.json file.
//
// 1. Start with a hardcoded file & tag list. (The same files from the build step in app/Makefile.)
// 2. Look for the those files in the data dir, and compute the md5 of each.
// 3. Write the result out to appData.json.
// 
// This is called during the build step of app/Makefile.

const fs = require("fs");
const child_process = require("child_process");
const { release } = require("os");

const dataDir = "../server/server_data"
const appData = {
    "releases": [
        {
            "version": "1.3",
            "tags": ["active", "odd"],
            "platform": "windows",
            "md5": "",
            "file": "app_odd.exe"
        },
        {
            "version": "1.2",
            "tags": ["archive", "even"],
            "platform": "windows",
            "md5": "",
            "file": "app_even.exe"
        },
        {
            "version": "1.3",
            "tags": ["active", "odd"],
            "platform": "darwin",
            "md5": "",
            "file": "app_odd_mac"
        },
        {
            "version": "1.2",
            "tags": ["lts", "even"],
            "platform": "darwin",
            "md5": "",
            "file": "app_even_mac"
        },
        {
            "version": "1.3",
            "tags": ["active", "odd"],
            "platform": "linux",
            "md5": "",
            "file": "app_odd_linux"
        },
        {
            "version": "1.2",
            "tags": ["lts", "even"],
            "platform": "linux",
            "md5": "",
            "file": "app_even_linux"
        }
    ]
};

for (const release of appData.releases) {
    const spawn = child_process.spawnSync("md5", ["-q", `${dataDir}/${release.file}`]);
    const errorText = spawn.stderr.toString().trim();
    const md5 = spawn.stdout.toString().trim();
    if (errorText){
        console.error(`Skipped getting md5 of ${release.file}: ${errorText}`)
    } else if (md5 !== ""){
        release.md5 = md5;
    }
}

const appDataFile = `${dataDir}/appData.json`;
try {
    fs.writeFileSync(appDataFile, JSON.stringify(appData, null, 2));
    console.log(`Created ${appDataFile}`);
}
catch (e) {
    console.error(`Error writing data file '${appDataFile}`, e);
}