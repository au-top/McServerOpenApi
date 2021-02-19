import express, { Express, Request, Response } from "express";
import path from "path";
import fs, { fstatSync } from "fs";
interface resAsync<T> {
    err: any;
    res: T;
}
async function modInfo(
    modsPath: string,
    clientModsPath: string,
    req: Request,
    res: Response
) {
    const resBody: {
        filename: string;
        length: number;
    }[] = [];
    //modsPath
    {
        const readRes = await new Promise<resAsync<string[]>>((res) =>
            fs.readdir(modsPath, (err, files) => res({ err, res: files }))
        );
        if (readRes.err != null) {
            console.log(readRes.err);
            return;
        }
        resBody.push(
            ...readRes.res
                .map((v) => path.join(modsPath, v))
                .filter((v) => {
                    try {
                        return fs.statSync(v).isFile();
                    } catch (e) {
                        console.log(e);
                        return false;
                    }
                })
                .map((v) => {
                    return { filename: path.parse(v).base, length: fs.statSync(v).size };
                })
        );
    }
    //clientModsPath
    {
        const readRes = await new Promise<resAsync<string[]>>((res) =>
            fs.readdir(clientModsPath, (err, files) => res({ err, res: files }))
        );
        if (readRes.err != null) {
            console.log(readRes.err);
            return;
        }
        resBody.push(
            ...readRes.res
                .map((v) => path.join(clientModsPath, v))
                .filter((v) => testIsFile(v))
                .map((v) => {

                    return { filename: path.parse(v).base, length: fs.statSync(v).size };
                })
        );
    }
    res.json(resBody);
}

function testIsFile(filePath: string) {
    try {
        const _fp = path.resolve(filePath);
        return fs.statSync(_fp).isFile();
    } catch (e) {
        console.error(e);
        return false;
    }
}

function createDir(dirPath: string) {
    if (!fs.existsSync(dirPath)) fs.mkdirSync(dirPath);
}

function setup(app: Express, gameServicPath: string) {
    const serviceGameModsPath = path.resolve(gameServicPath, "mods");
    const serviceGameClientModsPath = path.resolve(
        gameServicPath,
        "clientMods"
    );
    //init auto mkdir
    createDir(serviceGameModsPath);
    createDir(serviceGameClientModsPath);

    app.get("/modInfo", modInfo.bind(null, serviceGameModsPath,serviceGameClientModsPath));

    app.use("/modSource",express.static(serviceGameClientModsPath), express.static(serviceGameModsPath));

    console.log("openModSetup Successed");
}

export { setup };
