/*
    loading this module starts long polling and xhr queries (if appropriate)
    pollXhr implements long polling
        if event received, launch xhr query
    Each custom element adds its query part to gqlQuery (type GqlQuery)
        use GraphQL aliases for each custom element for clarity, especially if arguments differ
    GqlQuery dedups query parts and builds the xhr request body
    Upon successful xhr query, "refresh" event dispatched
        each custom element listens to refresh and renders itself when received
    event "cia3Error" dispatched on errors, with detail being a string which the Error custom element will add to itself
*/

let xhr = new XMLHttpRequest();
let pollXhr = new XMLHttpRequest();
let pollSince = Date.now() - 86400000;
const longPollTimeout = 30;
let data = {};

xhr.onload = () => {
	if (xhr.status >= 200 && xhr.status < 300) {
        data = JSON.parse(xhr.responseText).data;
        const refreshData = new CustomEvent("refresh");
        dispatchEvent(refreshData);
	} else {
        console.log(xhr);
        let cia3Error = new CustomEvent("cia3Error", { 'detail' : `Received non-2xx response on data query. Live updates will continue, but the latest save file is not shown here.`});
        dispatchEvent(cia3Error);
    }
}

xhr.onerror = e => {
    console.log(e);
    let cia3Error = new CustomEvent("cia3Error", { 'detail' : `Data query request failed. Live updates will continue, but the latest save file is not shown here.`});
    dispatchEvent(cia3Error);
}

let pollNow = () => {
    pollXhr.open('GET', `/events?timeout=${longPollTimeout}&category=refresh&since_time=${pollSince}`);
    pollXhr.send();
}

pollXhr.onload = () => {
    if (pollXhr.status >= 200 & pollXhr.status < 300) {
        let pollData = JSON.parse(pollXhr.responseText);
        if (typeof pollData.events != 'undefined') {
            pollSince = pollData.events[0].timestamp;
            xhr.open('POST', '/graphql');
            xhr.setRequestHeader('Content-Type', 'application/json');
            xhr.send(JSON.stringify(gqlQuery.body()));
        }
        if (pollData.timeout != undefined) {
            pollSince = pollData.timestamp;
        }
        pollNow();
    } else {
        console.log("failed xhr request:", pollXhr);
        let cia3Error = new CustomEvent("cia3Error", { 'detail' : `Received non-2xx response on polling query. Live updates have stopped.`});
        dispatchEvent(cia3Error);
    }
}

pollXhr.onerror = e => {
    console.error("Long poll returned error");
    console.log(e);
    let cia3Error = new CustomEvent("cia3Error", { 'detail' : `Polling error. Live updates have stopped. Correct and refresh page.`});
    dispatchEvent(cia3Error);
}

class GqlQuery {
    // Using Set to deduplicate queries
    queryParts = new Set();
    query() {
        return '{' + Array.from(this.queryParts).join('\n') + '}';
    }
    body() {
        return {
            'query' : this.query()
        }
    }
}
let gqlQuery = new GqlQuery();

// Most of the cia3-* elements follow this form, so extend this class
class Cia3Element extends HTMLElement {
    connectedCallback() {
        this.registerAndListen();
    }
    render() {
        this.innerText = 'REPLACE ME';
    }
    registerAndListen() {
        gqlQuery.queryParts.add(this.queryPart);
        window.addEventListener('refresh', () => this.render());
    }
    queryPart = 'REPLACE ME';
}

// TODO: Allow removal of error messages
class Error extends HTMLElement {
    connectedCallback() {
        window.addEventListener('cia3Error', (e) => this.render(e.detail));
    }
    render(errMsg) {
        const p = document.createElement('p');
        p.innerText = errMsg;
        this.appendChild(p);
    }
}

class Filename extends Cia3Element {
    render() {
        this.innerText = data.fileName;
    }
    queryPart = 'fileName';
}

class Fullpath extends Cia3Element {
    render() {
        this.innerText = data.fullPath;
    }
    queryPart = 'fullPath';
}

class Difficulty extends Cia3Element {
    render() {
        this.innerText = data.difficultyNames[data.difficulty[0]].name;
    }
    queryPart = `
        difficulty: int32s(section: "GAME", nth: 2, offset: 20, count: 1)
        difficultyNames: listSection(target: "bic", section: "DIFF", nth: 1) { name: string(offset:0, maxLength: 64) }
    `;
}

class Map extends Cia3Element {
    connectedCallback() {
        let spoilerMask = 0x2;
        this.queryPart = `
            map(playerSpoilerMask: ${spoilerMask}) {
                tileSetWidth
                tileSetHeight
                tiles {
                    hexTerrain
                    chopped
                }
            }
        `;
        this.registerAndListen();
    }
    render() {
        this.innerHTML = '';
        let tilesWide = Math.floor(data.map.tileSetWidth / 2);
        this.style.setProperty('--map-width', tilesWide);
        data.map.tiles.forEach( (e, i) => {
            const tile = document.createElement('cia3-tile');
            if (e.hexTerrain) {
                tile.setAttribute('data-terrain', e.hexTerrain);
            }
            if (e.chopped) {
                tile.setAttribute('data-chopped', 'true');
            }
            if ((i + tilesWide) % data.map.tileSetWidth == 0) {
                tile.classList.add('odd-row');
            }
            this.appendChild(tile);
        });
    }
    queryPart = '';
}

class Tile extends HTMLElement {
	connectedCallback () {
		this.render();
	}
    baseTerrainCss = {
        '0': 'desert',
        '1': 'plains',
        '2': 'grassland',
        '3': 'tundra',
        'b': 'coast',
        'c': 'sea',
        'd': 'ocean'
    }
    overlayTerrain = {
        '4': 'fp',
        '5': 'hill',
        '6': '⛰️',
        '7': '🌲',
        '8': '🌴',
        '9': 'marsh',
        'a': '🌋'
    }
    render () {
        const tileDiv = document.createElement('div');
        this.appendChild(tileDiv);
        tileDiv.classList.add('isotile');
        if (this.dataset.chopped == 'true') {
            const chopDiv = document.createElement('div');
            chopDiv.classList.add('chopped');
            this.appendChild(chopDiv);
        }
        let terr = this.dataset.terrain;
        if (terr) {
            if (this.baseTerrainCss[terr[1]]) {
                this.style.setProperty('--tile-color', `var(--${this.baseTerrainCss[terr[1]]})`);
            }
            if (this.overlayTerrain[terr[0]]) {
                const terrOverlayDiv = document.createElement('div');
                this.appendChild(terrOverlayDiv);
                terrOverlayDiv.className = 'terrain-overlay';
                terrOverlayDiv.innerText = this.overlayTerrain[terr[0]];
            }
        }
        let text = this.dataset.text;
        if (text) {
            const textDiv = document.createElement('div');
            textDiv.classList.add('tiletext');
            this.appendChild(textDiv);
        }
    }
}

class Url extends HTMLElement {
    connectedCallback() {
        this.render();
    }
    render() {
        let url = location.protocol + "//" + location.host;
        this.innerHTML = `<a href="${url}" target="_blank">${url}</a>`;
    }
}

// TODO: Add controls to customize query and re-query. And remove old query from gqlQuery.
class HexDump extends Cia3Element {
    render() {
        this.innerText = 'Hex dump tool under construction, no controls yet.\n' + data.cia3Hexdump;
    }
    queryPart = 'cia3Hexdump: hexDump(section: "CIV3", nth: 1, offset: -4, count: 2048)';
}

class MapX extends Cia3Element {
    render() {
        this.innerText = data.mapx[0];
    }
    queryPart = 'mapx: int32s(section: "WRLD", nth: 2, offset: 8, count: 1)';
}

class MapY extends Cia3Element {
    render() {
        this.innerText = data.mapy[0];
    }
    queryPart = 'mapy: int32s(section: "WRLD", nth: 2, offset: 28, count: 1)';
}

class WorldSize extends Cia3Element {
    render() {
        this.innerText = data.worldSizeNames[data.worldsize.size].name;
    }
    queryPart = `
        worldsize: civ3 { size }
        worldSizeNames: listSection(target: "bic", section: "WSIZ", nth: 1) { name: string(offset:32, maxLength: 32) }
    `;
}

class Barbarians extends Cia3Element {
    render() {
        this.innerText = this.barbariansSettings[data.barbarians.barbariansFinal.toString()];
    }
    queryPart = 'barbarians: civ3 { barbariansFinal }';
    barbariansSettings = {
        '-1': 'No Barbarians',
        '0': 'Sedentary',
        '1': 'Roaming',
        '2': 'Restless',
        '3': 'Raging',
        '4': 'Random'
    };
}

class WorldSeed extends Cia3Element {
    render() {
        this.innerText = data.worldseed.worldSeed;
    }
    queryPart = 'worldseed: civ3 { worldSeed }';
}

class LandMass extends Cia3Element {
    render() {
        this.innerText = this.landMassNames[data.landmass.landMassFinal];
    }
    queryPart = 'landmass: civ3 { landMassFinal }';
    landMassNames = [
        "Archipelago",
        "Continents",
        "Pangea",
        "Random"
    ];
}

class OceanCoverage extends Cia3Element {
    render() {
        this.innerText = this.oceanCoverageNames[data.oceancoverage.oceanCoverageFinal];
    }
    queryPart = 'oceancoverage: civ3 { oceanCoverageFinal }';
    oceanCoverageNames = [
        "80% Water",
        "70% Water",
        "60% Water",
        "Random"
    ];
}

class Climate extends Cia3Element {
    render() {
        this.innerText = this.climateNames[data.climate.climateFinal];
    }
    queryPart = 'climate: civ3 { climateFinal }';
    climateNames = [
        "Arid",
        "Normal",
        "Wet",
        "Random"
    ];
}

class Temperature extends Cia3Element {
    render() {
        this.innerText = this.temperatureNames[data.temperature.temperatureFinal];
    }
    queryPart = 'temperature: civ3 { temperatureFinal }';
    temperatureNames = [
        "Warm",
        "Temperate",
        "Cool",
        "Random"
    ];
}

class Age extends Cia3Element {
    render() {
        this.innerText = this.ageNames[data.age.ageFinal];
    }
    queryPart = 'age: civ3 { ageFinal }';
    ageNames = [
        "3 Billion",
        "4 Billion",
        "5 Billion",
        "Random"
    ];
}

class Civs extends Cia3Element {
    numFields = 64;
    render() {
        // this.innerHTML = JSON.stringify(data.civs, null, '  ');
        this.innerHTML = '';
        const table = document.createElement('table');
        const friendlyTable = document.createElement('table');
        const hexDumps = document.createElement('div');
        hexDumps.classList += "dump";
        this.appendChild(friendlyTable);
        this.appendChild(table);
        // table.innerHTML = '<tr><th>Player #</th><th>RACE ID</th>' + '<th>?</th>'.repeat(this.numFields - 2) + '</tr>';
        let headers = "";
        for (let i = 2; i < this.numFields; i++) {
            headers += `<th>${i} 0x${(i*4).toString(16)} ${i*4}</th>`
        }
        table.innerHTML = '<tr><th>Player #</th><th>RACE ID</th>' +  headers + '</tr>';
        friendlyTable.innerHTML = `<tr>
            <th>Player #</th>
            <th>Civ Name</th>
            <th>Government</th>
            <th>Mobilization</th>
            <th>Tiles Discovered</th>
            <th>Era</th>
        </tr>`;
        data.civs.filter(this.civsFilter).forEach((e, i) => {
            const friendlyRow = document.createElement('tr');
            friendlyRow.innerHTML += `<td>${e.playerNumber[0]}</td>`;
            friendlyRow.innerHTML += `<td>${data.race[e.raceId[0]].civName}</td>`;
            friendlyRow.innerHTML += `<td>${data.governmentNames[e.governmentType[0]].name}</td>`;
            friendlyRow.innerHTML += `<td>${e.mobilizationLevel[0]}</td>`;
            friendlyRow.innerHTML += `<td>${e.tilesDiscovered[0]}</td>`;
            friendlyRow.innerHTML += `<td>${data.eraNames[e.era[0]].name}</td>`;
            // friendlyRow.innerHTML += `<td>${}</td>`;
            // friendlyRow.innerHTML += `<td>${}</td>`;
            friendlyTable.appendChild(friendlyRow);
            const row = document.createElement('tr');
            e.int32s.forEach((ee, ii) => {
                const td = document.createElement('td');
                switch (ii) {
                    case 1:
                        if (ee != -1) {
                            td.innerText = data.race[ee].civName;
                            break;
                        }
                    default:
                        td.innerText = ee;
                }
                row.appendChild(td);
            });
            table.appendChild(row);
            if (this.oldCivsData != undefined) {
                const hexDiff = document.createElement('div');
                // let foo = Diff.diffWordsWithSpace(this.oldCivsData[e.int32s[0]].hexDump, data.civs[e.int32s[0]].hexDump);
                let foo = Diff.createTwoFilesPatch("old", "new" ,this.oldCivsData[e.int32s[0]].hexDump, data.civs[e.int32s[0]].hexDump);
                // foo.forEach(function(part){
                //     // green for additions, red for deletions
                //     // grey for common parts
                //     let color = part.added ? 'green' :
                //       part.removed ? 'red' : 'grey';
                //     let span = document.createElement('span');
                //     span.style.color = color;
                //     span.appendChild(document
                //       .createTextNode(part.value));
                //     hexDiff.appendChild(span);
                // });
                let diff2Html = Diff2Html.html(foo, {
                    drawFileList: true,
                    matching: 'words',
                    outputFormat: 'side-by-side',
                });
                hexDiff.innerHTML = '*****\n#' + e.int32s[0] + ' : ' + data.race[e.int32s[1]].civName + ' Diff\n*****\n' + diff2Html;
                this.appendChild(hexDiff);
            }
            hexDumps.innerHTML += '*****\n#' + e.int32s[0] + ' : ' + data.race[e.int32s[1]].civName + '\n*****\n' +
                e.hexDump.replace(/((?: 00)+)/g, '<span class="dim">$1</span>')
                .replace(/(\.+)/g, '<span class="dim">$1</span>')
                .replace(/((?: ff)+)/g, '<span class="medium">$1</span>')
                ;
        })
        this.appendChild(hexDumps);
        this.oldCivsData = data.civs
    }
    civsFilter (e) {
        return e.int32s[1] > 0; // non-barb players
        // return e.int32s[0] == 1; // first human player only
    }
    oldCivsData = undefined;
    queryPart = `
        civs {
            int32s(offset:0, count: ${this.numFields})
            hexDump
            playerNumber: int32s(offset:0, count: 1)
            raceId: int32s(offset:4, count: 1)
            governmentType: int32s(offset:132, count: 1)
            mobilizationLevel: int32s(offset:136, count: 1)
            tilesDiscovered: int32s(offset:140, count: 1)
            era: int32s(offset:252, count: 1)
            UNSUREresearchBulbs: int32s(offset:256, count: 1)
            UNSUREcurrentResearchId: int32s(offset:260, count: 1)
            UNSUREcurrentResearchTurns: int32s(offset:264, count: 1)
            UNSUREfutureTechsCount: int32s(offset:268, count: 1)
        }
        race {
            leaderName: string(offset:0, maxLength: 32)
            leaderTitle: string(offset:32, maxLength: 24)
            adjective:  string(offset:88, maxLength: 40)
            civName: string(offset:128, maxLength: 40)
            objectNoun: string(offset:168, maxLength: 40)
        }
        governmentNames: listSection(target: "bic", section: "GOVT", nth: 1) { name: string(offset:24, maxLength: 64) }
        eraNames: listSection(target: "bic", section: "ERAS", nth: 1) { name: string(offset:0, maxLength: 64) }
    `;
}

window.customElements.define('cia3-error', Error);
window.customElements.define('cia3-filename', Filename);
window.customElements.define('cia3-fullpath', Fullpath);
window.customElements.define('cia3-difficulty', Difficulty);
window.customElements.define('cia3-map', Map);
window.customElements.define('cia3-tile', Tile);
window.customElements.define('cia3-url', Url);
window.customElements.define('cia3-hexdump', HexDump);
window.customElements.define('cia3-mapx', MapX);
window.customElements.define('cia3-mapy', MapY);
window.customElements.define('cia3-worldsize', WorldSize);
window.customElements.define('cia3-barbarians', Barbarians);
window.customElements.define('cia3-worldseed', WorldSeed);
window.customElements.define('cia3-landmass', LandMass);
window.customElements.define('cia3-oceancoverage', OceanCoverage);
window.customElements.define('cia3-climate', Climate);
window.customElements.define('cia3-temperature', Temperature);
window.customElements.define('cia3-age', Age);
window.customElements.define('cia3-civs', Civs);
pollNow();
