(function (window) {
    const document = window.document;
    let goLink, errorNode;

    document.addEventListener('DOMContentLoaded',function () {
        goLink = document.querySelector("#go");
        errorNode = document.querySelector("#error");
    });

    function Inputs(names) {
        this._names = names || [];
        for (name of this._names) {
            this[name] = document.querySelector('#' + name);
        }
        console.log('inputs:',this);
    }

    Inputs.prototype.clear = function () {
        for (name of this._names) {
            if (this[name].type === 'checkbox') {
                this[name].checked = "";
            } else {
                this[name].value = "";
            }
        }
    };

    Inputs.prototype.params = function () {
        const params = {};
        for (name of this._names) {
            if (this[name].type === 'checkbox') {
                if (this[name].checked) {
                    params[name] = this[name].value;
                }
            } else {
                params[name] = this[name].value;
            }
        }
        return params;
    };

    Inputs.prototype.paramsURL = function (base) {
        const params = this.params();
        const qp = [];
        for (const name in params) {
            if (params.hasOwnProperty(name) && params[name] !== '') {
                qp.push(name + '=' + escape(params[name]));
            }
        }
        return base + '?' + qp.join('&');
    };

    const progress = {
        start: function () {
            errorNode.innerHTML = '';
            goLink.disabled = true;
            goLink.src = '/images/progress.svg?bust=1';
        },
        end: function () {
            goLink.disabled = false;
            goLink.src = '/images/go.svg?bust=1';
        }
    };

    const createDefinitionLink = function (word) {
        const document = window.document;
        if (word.includes(' ')) {
            return document.createTextNode(word);
        }
        const link = document.createElement('a');
        link.target = 'definition';
        link.href = 'https://www.thewordfinder.com/define/' + escape(word);
        link.innerText = word;
        return link;
    };

    const appendResults = function (parent, words, prop) {
        let i = 0;
        for (const entry of words) {
            i++;
            const word = prop ? entry[prop] : entry;
            const span = document.createElement("span");
            span.classList.add('inline-result');
            span.classList.add(i % 2 === 0 ? 'even': 'odd');
            span.appendChild(crossword.createDefinitionLink(word));
            parent.appendChild(span);
        }
    }

    const callAPI = function (url, renderFn) {
        progress.start();
        const xhr = new window.XMLHttpRequest();
        xhr.responseType = "json";
        xhr.onreadystatechange = function () {
            if (xhr.readyState !== 4) return;
            console.log('response for get', url);
            console.log(xhr.response);
            progress.end();
            if (xhr.status >= 200 && xhr.status < 300) {
                renderFn(xhr.response);
            } else {
                console.warn('GET',url,'failed');
                console.warn('response',xhr.response);
                errorNode.innerHTML = xhr.response.error;
            }
        };
        xhr.open('GET', url);
        xhr.send();
    };

    window.crossword = {
        progress: progress,
        createDefinitionLink: createDefinitionLink,
        appendResults: appendResults,
        callAPI: callAPI,
        Inputs: Inputs,
    };
})(window)