function showMenu() {
    document.getElementById("mobile-menu").setAttribute("aria-expanded", "true")
}

function hideMenu() {
    document.getElementById("mobile-menu").setAttribute("aria-expanded", "false")
}

function setDirty() {
    document.getElementById("save-button").classList.remove("invisible")
}

function setClean() {
    document.getElementById("save-button").classList.add("invisible")
}