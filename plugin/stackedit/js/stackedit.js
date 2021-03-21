var div = document.getElementById("id_content");
if (div) {
    var el = document.createElement("span");
    el.classList.add("helptext");
    el.innerHTML = '<button type="button" onclick="openStackEditor(\'id_content\');">Markdown editor</button> |';
    insertAfter(el, div);
}

function insertAfter(newNode, referenceNode) {
    referenceNode.parentNode.insertBefore(newNode, referenceNode.nextSibling);
}
function openStackEditor(elementId) {
    const el = document.querySelector(`textarea[id = '${elementId}']`);
    const stackedit = new Stackedit();

    // Open the iframe.
    stackedit.openFile({
        content: {
            text: el.value,
            properties: {
                extensions: {
                    preset: 'commonmark',
                    markdown: {
                        table: true,
                    }
                },
                colorTheme: 'dark',
            }
        }
    });

    // Listen to StackEdit events and apply the changes to the textarea.
    stackedit.on('fileChange', (file) => {
        el.value = file.content.text;
    });
}