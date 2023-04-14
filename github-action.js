const core = require("@actions/core");
const fs = require("fs").promises;
const path = require("path");
const http = require("http");
const https = require("https");
const tar = require("tar");
const { ok } = require("assert");

function download(url, cnt) {
  return new Promise((rs, rj) => {
    if (cnt > 8) {
      rj(Error(`loop detected ${url}`));
      return;
    }
    let client = http;
    if (url.toString().indexOf("https") === 0) {
      client = https;
    }
    client
      .request(url, (response) => {
	if (300<=response.statusCode && response.statusCode < 400) {
		const url = response.headers['location'];
		console.log(`url ==> ${url}`);
		rs(download(url, cnt + 1));
		return;
	}
	if (!(200<=response.statusCode && response.statusCode < 300)) {
		rj(`fetch error ${url}:${response.statusCode}`);
	}
        response.on("error", (e) => {
          rj(e);
        });
        const data = [];
        response.on("data", function (chunk) {
          data.push(Buffer.from(chunk));
        });
        response.on("end", function () {
          rs(Buffer.concat(data));
        });
      })
      .end();
  });
}

function getTempDirectory() {
  const tempDirectory = process.env["RUNNER_TEMP"] || ".";
  ok(tempDirectory, "Expected RUNNER_TEMP to be defined");
  return tempDirectory;
}

async function main() {
  try {
    const params = {
	    version: core.getInput("version"),
	    url: core.getInput("url"),
	    filename: core.getInput("filename"),
	    os: core.getInput("os"), // no default
	    suffix: core.getInput("suffix"),
	    cpu: core.getInput("cpu") // no default
    };
    if (params.os === "") {
	    switch (process.platform) {
	      case "darwin":
		params.os = "darwin";
		break;
	      case "linux":
		params.os = "linux";
	    }
     }
     if (params.cpu === "") {
	switch (process.arch) {
	  case "i386":
	    params.cpu = "i386";
	    break;
	  case "x64":
	    params.cpu = "amd64";
	    break;
	  case "arm":
	    params.cpu = "armv7";
	    break;
	  case "arm64":
	    params.cpu = "arm64";
	    break;
	}
    }
    core.info(`Fetch hcloud choosen cpu [${params.cpu}][${process.arch}]`)
    const plainVersion = params.version.replace(/^v/, '');
    const hcloudUrl = `${params.url}/${params.version}/${params.filename}-${params.os}-${params.cpu}${params.suffix}`;
    core.info(`Fetch hcloud from:[${hcloudUrl}]`)
    const hcloudBin = await download(hcloudUrl, 0);
    const hcloudBinDir = path.join(getTempDirectory(), "hcloud-bin");
    // const dir = path.join(getTempDirectory()); // , `hcloud${params.suffix}`);
    await fs.mkdir(hcloudBinDir, {
      recursive: true,
      mode: 0o755,
    });
    const hcloudFnameTar = path.join(getTempDirectory(), `hcloud${params.suffix}`);
    await fs.writeFile(hcloudFnameTar, hcloudBin);
    await tar.x({
	    file: hcloudFnameTar,
	    cwd: hcloudBinDir
    });
    // await fs.chmod(hcloudFname, 0o755);
    core.exportVariable("HCLOUD_URL", hcloudUrl);
    const hcloudFname = path.join(hcloudBinDir, "hcloud");
    core.exportVariable("HCLOUD_FNAME", hcloudFname);
    core.addPath(hcloudBinDir);
    await fs.unlink(hcloudFnameTar);
    core.info(`Installed hcloud into:[${hcloudFname}]`)
  } catch (e) {
    core.setFailed(e);
  }
}

main();
