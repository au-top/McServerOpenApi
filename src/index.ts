import express from "express";
import fs from "fs";
import path from "path";
const appDirPath = __dirname;

const pluginDirFileList = fs.readdirSync(`${appDirPath}/plugin`);
const port=25503;

const config=JSON.parse( fs.readFileSync(path.resolve(__dirname,"./config.json")).toString() ) as config;



const readJsModuleFile = pluginDirFileList
    .map((v) => `${appDirPath}/plugin/${v}`)
    .filter((v) => {
        const rfstate = fs.statSync(v);
        return rfstate.isFile() && /[\s\S]{1,}\.(js)|(ts)$/.test(v);
    });

(async () => {
    const app = express();
    const loadPluginMapList=await Promise.all(
        readJsModuleFile.map((v) => {
            return (async (v) => {
                const loadPlugin = await import(v);
                if (Object.prototype.hasOwnProperty.call(loadPlugin, "setup")) {
                    loadPlugin["setup"](app,config["mcServerPath"]);
                    return {
                        'pluginPath':v,
                        'loadState':'success'
                    }
                } else {
                    console.log(v, "setup error", "no find setup func");
                    return {
                        'pluginPath':v,
                        'loadState':'error'
                    }
                }
            })(v);
        })
    );
    console.log('load plugin success');
    console.log('load plugin list \n',loadPluginMapList);
    app.listen(25503,'0.0.0.0',()=>console.log(
        `server run listen port ${port}`
    ));
})().catch(console.log);
