const fallbackWidth = 80;
const minWidth = 20;
const maxWidth = 300;
const sampleText = "MMMMMMMMMM";
const themeStorageKey = "ascii-art-theme";
const themeChannelName = "ascii-art-theme-sync";

function isSupportedTheme(theme) {
    return theme === "dark" || theme === "light";
}

function readStoredTheme() {
    try {
        const storedTheme = window.sessionStorage.getItem(themeStorageKey);

        if (isSupportedTheme(storedTheme)) {
            return storedTheme;
        }
    } catch {
        // Storage may be unavailable; light remains the safe default.
    }
    return "light";
}

function storeTheme(theme) {
    try {
        window.sessionStorage.setItem(themeStorageKey, theme);
    } catch {
        // The current page can still change theme without persistence.
    }
}

// Keep the CSS theme, accessible state, and visible status synchronized.
function applyTheme(theme, themeToggle) {
    const themeStatus = themeToggle.querySelector(
        "[data-theme-toggle-status]"
    );

    if (theme === "dark") {
        document.documentElement.dataset.theme = "dark";
        themeToggle.setAttribute("aria-pressed", "true");

        if (themeStatus) {
            themeStatus.textContent = "ON";
        }
    } else {
        delete document.documentElement.dataset.theme;
        themeToggle.setAttribute("aria-pressed", "false");

        if (themeStatus) {
            themeStatus.textContent = "OFF";
        }
    }
}

document.addEventListener("DOMContentLoaded", function () {
    // Width measurement and theme setup are independent; missing markup for
    // one feature must not disable the other.
    const form = document.querySelector(".generation-form");
    const widthInput = document.querySelector('input[name="width"]');
    const terminal = document.querySelector(".terminal-output");
    const themeToggle = document.querySelector("[data-theme-toggle]");

    if (form && widthInput && terminal) {
        if (!widthInput.value) {
            widthInput.value = String(fallbackWidth);
        }

        form.addEventListener("submit", function () {
            widthInput.value = String(measureColumns(terminal));
        });
    }

    if (themeToggle) {
        let currentTheme = readStoredTheme();
        applyTheme(currentTheme, themeToggle);

        // sessionStorage is tab-scoped; BroadcastChannel shares the active
        // preference only with other tabs that are currently open.
        let themeChannel = null;
        if ("BroadcastChannel" in window) {
            themeChannel = new BroadcastChannel(themeChannelName);

            themeChannel.addEventListener("message", function (event) {
                const message = event.data;

                if (!message || typeof message !== "object") {
                    return;
                }

                if (message.type === "theme-request") {
                    themeChannel.postMessage({
                        type: "theme-update",
                        theme: currentTheme,
                    });
                    return;
                }

                if (
                    message.type === "theme-update" &&
                    isSupportedTheme(message.theme)
                ) {
                    currentTheme = message.theme;
                    applyTheme(currentTheme, themeToggle);
                    storeTheme(currentTheme);
                }
            });

            // With no open peer to answer, this tab keeps its local theme.
            themeChannel.postMessage({
                type: "theme-request",
            });
        }

        themeToggle.addEventListener("click", function () {
            let nextTheme = "dark";

            if (currentTheme === "dark") {
                nextTheme = "light";
            }

            currentTheme = nextTheme;
            applyTheme(currentTheme, themeToggle);
            storeTheme(currentTheme);

            if (themeChannel) {
                themeChannel.postMessage({
                    type: "theme-update",
                    theme: currentTheme,
                });
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
