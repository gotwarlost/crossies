(function (window) {
    const document = window.document;
    const crossword = window.crossword;
    const apiPath = '/api/v1/matching-words';
    const names = [ 'frame', 'page', 'syn1', 'syn2' ];

    let controls = {};
    let inputs;

    function showTiles(parent) {
        const frame = inputs.frame.value;
        for (let i=0; i < frame.length; i++) {
            const span = document.createElement('span');
            span.classList.add('tile');
            const ch = frame.substring(i,i+1);
            if (ch === '.') {
                span.innerHTML = '&nbsp;'
            } else {
                span.innerText = ch;
            }
            parent.appendChild(span);
        }
        if (frame.length > 0) {
            const span = document.createElement('span');
            span.classList.add('tile');
            span.classList.add('count');
            span.appendChild(document.createTextNode('(' + frame.length + ')'));
            parent.appendChild(span);
        }
    }

    function addSummaryText(result) {
        const things = result.totalWords === 1 ? 'word' : 'words';
        controls.resultsTitle.innerHTML = '';
        const suffix = document.createTextNode(' ' + result.totalWords + ' ' + things + ' found');
        showTiles(controls.resultsTitle);
        controls.resultsTitle.appendChild(document.createElement('br'));
        controls.resultsTitle.appendChild(suffix);
    }

    function renderResults(result) {
        controls.queryPage.classList.add('hide');
        controls.cleanButton.classList.add('hide');
        controls.resultsPage.classList.remove('hide');
        controls.backLink.classList.remove('hide');

        addSummaryText(result);

        if (result.query.synonyms && result.query.synonyms.length > 0) {
            const title = document.createElement('div');
            title.classList.add('row');
            controls.resultsNode.appendChild(title);
            if (result.synonymMatches && result.synonymMatches.length > 0) {
                const matches = result.synonymMatches;
                const things2 = matches.length === 1 ? 'word' : 'words';
                title.appendChild(document.createTextNode(' ' + matches.length + ' synonymous ' + things2 + ' found'));
                const synResults = document.createElement('div');
                synResults.classList.add('row');
                controls.resultsNode.appendChild(synResults);
                crossword.appendResults(synResults, matches);
            } else {
                title.appendChild(document.createTextNode('no words matched synonyms'));
            }
            const title2 = document.createElement('div');
            title2.classList.add('row');
            title2.appendChild(document.createTextNode('all matches'));
            controls.resultsNode.appendChild(title2);
        }

        const div = document.createElement("div");
        div.classList.add('row');
        controls.resultsNode.appendChild(div);
        crossword.appendResults(div, result.words);
        inputs.page.value = String(result.nextPage);
        if (result.nextPage !== 0) {
            controls.moreLink.classList.remove('hide');
        } else {
            controls.moreLink.classList.add('hide');
        }
    }

    function displayInputs() {
        controls.showWord.innerHTML = '';
        showTiles(controls.showWord);
    }

    document.addEventListener('DOMContentLoaded', function () {
        controls = {
            form: document.querySelector("#query"),
            resultsNode: document.querySelector('#results-list'),
            queryPage: document.querySelector('#query-page'),
            resultsPage: document.querySelector('#results-page'),
            backLink: document.querySelector('#back-link'),
            resultsTitle: document.querySelector('#results-title'),
            moreLink: document.querySelector('#more-link'),
            showWord: document.querySelector('#show-word'),
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
        const doQuery = function () {
            crossword.callAPI(inputs.paramsURL(apiPath), renderResults);
            return false;
        }
        controls.form.onsubmit = function () {
            controls.resultsNode.innerHTML = '';
            inputs.page.value = "1"
            return doQuery();
        };
        controls.moreLink.onclick = doQuery;
        inputs.frame.focus();
        inputs.frame.setSelectionRange(0, 0);
        inputs.frame.addEventListener('keyup', displayInputs);
        controls.cleanButton.addEventListener('click', function () { inputs.clear(); displayInputs(); });
    });
})(window)
