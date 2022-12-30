const fs = require("fs");


const allfiles = [
	"aes_helper.cl",
	"blake.cl",
	"bmw.cl",
	"cubehash.cl",
	"echo.cl",
	"fugue.cl",
	"groestl.cl",
	"hamsi.cl",
	"hamsi_help.cl",
	"hamsi_helper.cl",
	"hamsi_helper_big.cl",
	"jh.cl",
	"keccak.cl",
	"luffa.cl",
	"sha2_512.cl",
	"sha3_256.cl",
	"shabal.cl",
	"shavite.cl",
	"simd.cl",
	"skein.cl",
	"whirlpool.cl",
	"x16rs.cl",
	"x16rs_main.cl",
	"x16rs_main_empty_test.cl",
	];

let gopkgfilecon = `
package worker
// 输出所有 opencl 文件
func GetRenderCreateAllOpenclFiles() map[string]string {
	files := map[string]string{}

`;

// 依次读取文件
for (let i in allfiles) {
	let fname = allfiles[i];
	let fdata = fs.readFileSync(fname).toString();
	gopkgfilecon += `
	files["`+fname+`"] = \``+fdata+`\`
	`

}

gopkgfilecon += `
	return files
}
`;

// 写入文件
fs.writeFileSync("./worker/pkgcls.go", gopkgfilecon)





