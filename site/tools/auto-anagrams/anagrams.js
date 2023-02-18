(function (window) {
    const document = window.document;
    const crossword = window.crossword;
    const apiPath = '/api/v1/anagrams';
    const names = [ 'phrase', 'partial' ];
    let controls = {};
    let inputs;

    function renderResults(response) {
        controls.backLink.classList.remove('hide');
        controls.queryPage.classList.add('hide');
        controls.cleanButton.classList.add('hide');

        const entries = response.phrases;
        const things = entries.length === 1 ? 'anagram' : 'anagrams';
        controls.resultsTitle.innerHTML = '';
        controls.resultsTitle.appendChild(document.createTextNode('Anagrams for: ' + inputs.phrase.value));
        controls.resultsTitle.appendChild(document.createElement('br'));
        controls.resultsTitle.appendChild(document.createTextNode('' + entries.length + ' ' + things + ' found'));
        controls.resultsPage.classList.remove('hide');
        controls.resultsNode.innerHTML = '';
        const div = document.createElement("div");
        div.setAttribute("class","row");
        controls.resultsNode.appendChild(div);
        crossword.appendResults(div, entries);
    }

    document.addEventListener('DOMContentLoaded', function () {
        inputs = new crossword.Inputs(names);
        controls = {
            form: document.querySelector("#query"),
            resultsNode: document.querySelector('#results-list'),
            queryPage: document.querySelector('#query-page'),
            resultsPage: document.querySelector('#results-page'),
            backLink: document.querySelector('#back-link'),
            resultsTitle: document.querySelector('#results-title'),
            cleanButton: document.querySelector('#clean-link'),
        };
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
        inputs.phrase.focus();
        controls.cleanButton.addEventListener('click', function () { inputs.clear(); });
    });
})(window)
