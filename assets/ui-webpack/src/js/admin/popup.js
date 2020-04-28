/*!
 * HugoCMS
 * https://github.com/dfkdream/HugoCMS
 *
 * Copyright 2020 dfkdream
 * Released under MIT License
 */

require("spectre.css/dist/spectre-exp.css");
require("babel-polyfill");

async function uploadFiles(path,files,callback){
    for (let i=0;i<files.length;i++) {
        let file = files[i];
        callback(i+1,files.length,file.name);
        let resp = await fetch(path, {
            method: "POST",
            body: file
        });

        if (!resp.ok) {
            let data = await resp.json();
            alert(`Upload error: ${file.name} ${data.code} ${data.message}`);
        }
    }
}

module.exports = {
    confirm: (title, msg)=>{
        return new Promise((resolve)=>{
            let fragment = document.createDocumentFragment();
            let popup = document.createElement("div");
            popup.setAttribute("class","modal active");
            popup.innerHTML = require("../../html/confirm.html");
            fragment.append(popup);

            fragment.getElementById("title").innerText = title;
            fragment.getElementById("content").innerText = msg;

            let yes = fragment.getElementById("yes");
            yes.onclick = ()=>{
                document.body.removeChild(popup);
                resolve(true);
            };

            fragment.getElementById("no").onclick = ()=>{
                document.body.removeChild(popup);
                resolve(false);
            };

            fragment.getElementById("close-overlay").onclick = ()=>{
                document.body.removeChild(popup);
                resolve(false);
            };

            document.body.appendChild(fragment);

            yes.focus();
        });
    },

    prompt: (title,msg)=>{
        return new Promise((resolve)=>{
            let fragment = document.createDocumentFragment();
            let popup = document.createElement("div");
            popup.setAttribute("class","modal active");
            popup.innerHTML = require("../../html/prompt.html");
            fragment.append(popup);

            let input = fragment.getElementById("input");

            fragment.getElementById("title").innerText = title;
            fragment.getElementById("content").innerText = msg;

            let ok = fragment.getElementById("ok");
            ok.onclick = ((input)=>{
                return ()=>{
                    let v = input.value;
                    document.body.removeChild(popup);
                    resolve(v);
                };
            })(input);

            fragment.getElementById("cancel").onclick = ()=>{
                document.body.removeChild(popup);
                resolve(null);
            };

            fragment.getElementById("close-overlay").onclick = ()=>{
                document.body.removeChild(popup);
                resolve(null);
            };

            fragment.getElementById("form").onsubmit = ()=>{
                ok.onclick(null);
            };

            document.body.appendChild(fragment);

            input.focus();
        });
    },

    alert:(title,msg)=>{
        return new Promise((resolve)=>{
            let fragment = document.createDocumentFragment();
            let popup = document.createElement("div");
            popup.setAttribute("class","modal active");
            popup.innerHTML = require("../../html/alert.html");
            fragment.append(popup);

            fragment.getElementById("title").innerText = title;
            fragment.getElementById("content").innerText = msg;
            let ok = fragment.getElementById("ok");
            ok.onclick = ()=>{
                document.body.removeChild(popup);
                resolve();
            };

            fragment.getElementById("close-overlay").onclick = ()=>{
                document.body.removeChild(popup);
                resolve();
            };

            document.body.appendChild(fragment);

            ok.focus();
        });
    },

    upload:(uploadPath)=>{
        return new Promise((resolve)=>{
            let fragment = document.createDocumentFragment();
            let popup = document.createElement("div");
            popup.setAttribute("class","modal active");
            popup.innerHTML = require("../../html/upload.html");
            fragment.append(popup);

            const uploadModalTitle = fragment.getElementById("upload-modal-title");
            const uploadModalContent = fragment.getElementById("upload-modal-content");
            const uploadModalUpload = fragment.getElementById("upload-modal-upload");
            const uploadModalClose = fragment.getElementById("upload-modal-close");

            uploadModalUpload.setAttribute("disabled","");
            uploadModalClose.removeAttribute("disabled");
            uploadModalTitle.innerText = "Upload File";
            uploadModalContent.innerHTML="";
            let inputFile = document.createElement("input");
            inputFile.setAttribute("type","file");
            inputFile.setAttribute("class","form-input");
            inputFile.setAttribute("multiple","");
            uploadModalContent.appendChild(inputFile);

            inputFile.onchange = ()=>{
                uploadModalUpload.removeAttribute("disabled");
            };

            uploadModalUpload.onclick = ()=>{
                uploadModalTitle.innerText = "Uploading...";
                uploadModalContent.innerHTML="";
                uploadModalClose.setAttribute("disabled","");
                uploadModalUpload.setAttribute("disabled","");
                let statusProgress = document.createElement("progress");
                statusProgress.setAttribute("class","progress");
                statusProgress.setAttribute("max","100");
                uploadModalContent.appendChild(statusProgress);

                let statusP = document.createElement("p");
                uploadModalContent.appendChild(statusP);

                uploadFiles(uploadPath,inputFile.files,(idx,len,name)=>{
                    statusProgress.value = idx/len*100;
                    statusP.innerText = `(${idx}/${len}) ${name}`;
                    statusP.style.margin=".8rem 0 0";
                })
                    .then(()=>{
                        document.body.removeChild(popup);
                        resolve();
                    });
            };

            uploadModalClose.onclick = ()=>{
                document.body.removeChild(popup);
                resolve();
            };

            document.body.appendChild(fragment);
        })
    }
};