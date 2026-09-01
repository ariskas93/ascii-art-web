const fallbackWidth = 80;
const minWidth = 20;
const maxWidth = 300;
const sampleText = "MMMMMMMMMM";

document.addEventListener("DOMContentLoaded", function () {
    // Resolve the few elements C owns once. Missing elements simply leave the
    // server-side 80-column fallback in place.
    const form = document.querySelector(".generation-form");
    const widthInput = document.querySelector('input[name="width"]');
    const terminal = document.querySelector(".terminal-output");
    const themeToggle = document.querySelector("[data-theme-toggle]");

    if (!form || !widthInput || !terminal) {
        return;
    }
    if (!widthInput.value) {
        widthInput.value = String(fallbackWidth);
    }
    form.addEventListener("submit", function () {
        widthInput.value = String(measureColumns(terminal));
    });
    if (themeToggle) {
        themeToggle.addEventListener("click", function () {
            // The CSS listens to html[data-theme="dark"], so the button only
            // flips that attribute and its own accessible pressed state.
            const isDark = document.documentElement.dataset.theme !== "dark";

            if (isDark) {
                document.documentElement.dataset.theme = "dark";
                themeToggle.setAttribute("aria-pressed", "true");
                themeToggle.textContent = "Light theme";
            } else {
                delete document.documentElement.dataset.theme;
                themeToggle.setAttribute("aria-pressed", "false");
                themeToggle.textContent = "Dark theme";
            }
        });
    }
});

function measureColumns(terminal) {
    // Measure at submit time only; resizing later does not trigger requests or
    // realign an already-rendered snapshot.
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
    // A hidden probe uses the terminal font so pixel width can be converted to
    // text columns without rendering or aligning any ASCII art.
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
