import express, { Express, Request, Response } from "express";
import path from "path";
import fs, { fstatSync } from "fs";
interface resAsync<T> {
    err: any;
    res: T;
}
async function modInfo(modsPath: string, req: Request, res: Response) {
    const readRes = await new Promise<resAsync<string[]>>((res) =>
        fs.readdir(modsPath, (err, files) => res({ err, res:files }))
    );
    if (readRes.err != null) {
        console.log(readRes.err);
        return;
    }
    const resBody=readRes.res.filter((v) => {
        v=path.resolve(modsPath,v);
        try{
            if (fs.statSync(v).isFile()) return true;
            else return false;
        }catch(e){
            console.log(e);
            return false;
        }
    }).map(v=>{
        const vpath=path.resolve(modsPath,v);
        return {'filename':v,'length':fs.statSync(vpath).size}
    });
    res.json(resBody);
}

function setup(app: Express, gameServicPath: string) {
    const serviceGameModsPath = path.resolve(gameServicPath, "mods");
 
    app.get("/modInfo", modInfo.bind(null, serviceGameModsPath));

    app.use("/modSource", express.static(serviceGameModsPath));

    console.log("openModSetup Successed");
}

export { setup };
