(function (window) {
    const document = window.document;
    const crossword = window.crossword;
    const apiPath = '/api/v1/synonyms';
    const names = [ 'word', 'sort', 'minLetters', 'maxLetters', 'pattern', 'startsWith', 'endsWith', 'all' ];

    let controls = {};
    let inputs;

    function renderResults(response) {
        controls.queryPage.classList.add('hide');
        controls.backLink.classList.remove('hide');
        controls.cleanButton.classList.add('hide');

        const entries = response.entries || [];
        const things = entries.length === 1 ? 'entry' : 'entries';
        controls.resultsTitle.innerHTML = '';
        controls.resultsTitle.appendChild(document.createTextNode('Synonyms for: ' + inputs.word.value));
        controls.resultsTitle.appendChild(document.createElement('br'));
        controls.resultsTitle.appendChild(document.createTextNode('' + entries.length + ' ' + things + ' found'));
        controls.resultsPage.classList.remove('hide');
        controls.resultsNode.innerHTML = '';
        const div = document.createElement("div");
        div.setAttribute("class","row");
        controls.resultsNode.appendChild(div);
        crossword.appendResults(div, entries, 'synonym');
    }

    document.addEventListener('DOMContentLoaded', function () {
        controls = {
            form: document.querySelector("#query"),
            resultsNode: document.querySelector('#results-list'),
            queryPage: document.querySelector('#query-page'),
            resultsPage: document.querySelector('#results-page'),
            backLink: document.querySelector('#back-link'),
            resultsTitle: document.querySelector('#results-title'),
            cleanButton: document.querySelector('#clean-link'),
        };
        inputs = new crossword.Inputs(names);
        controls.resultsPage.classList.add('hide');
        controls.backLink.onclick = function () {
            controls.queryPage.classList.remove('hide');
            controls.resultsPage.classList.add('hide');
            controls.resultsNode.innerHTML = '';
            controls.backLink.classList.add('hide');
            controls.cleanButton.classList.remove('hide');
        }
        controls.form.onsubmit = function () {
            crossword.callAPI(inputs.paramsURL(apiPath), renderResults);
            return false;
        }
        inputs.word.focus();
        controls.cleanButton.addEventListener('click', function () { inputs.clear(); });
    });
})(window)
