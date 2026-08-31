const fallbackWidth = 80;
const minWidth = 20;
const maxWidth = 300;
const sampleText = "MMMMMMMMMM";

document.addEventListener("DOMContentLoaded", function () {
    const form = document.querySelector(".generation-form");
    const widthInput = document.querySelector('input[name="width"]');
    const terminal = document.querySelector(".terminal-output");

    if (!form || !widthInput || !terminal) {
        return;
    }
    if (!widthInput.value) {
        widthInput.value = String(fallbackWidth);
    }
    form.addEventListener("submit", function () {
        widthInput.value = String(measureColumns(terminal));
    });
});

function measureColumns(terminal) {
    const containerWidth = terminal.clientWidth;
    const charWidth = measureCharacterWidth(terminal);
    
    if (containerWidth <= 0 || charWidth <= 0) {
        return fallbackWidth;
    }
    const columns = Math.floor(containerWidth / charWidth);
    return clampWidth(columns);
}

function clampWidth(columns) {
    if (!Number.isFinite(columns)) {
        return fallbackWidth;
    }
    if (columns < minWidth) {
        return minWidth;
    }
    if (columns > maxWidth) {
        return maxWidth;
    }
    return columns;
}

function measureCharacterWidth(terminal) {
    const probe = document.createElement("span");
    probe.style.position = "absolute";
    probe.style.visibility = "hidden";
    probe.style.whiteSpace = "pre";
    probe.style.font = window.getComputedStyle(terminal).font;
    
    probe.textContent = sampleText;
    terminal.appendChild(probe);
    const width = probe.getBoundingClientRect().width / sampleText.length;
    probe.remove();
    if (!Number.isFinite(width)) {
        return 0;
    }
    return width;
}