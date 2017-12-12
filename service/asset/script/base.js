(function() {
    const e = document.querySelector.bind(document)

    // const log = () => []
    const log = console.log.bind(console)

    window.onload = function() {
        var s = document.createElement('p')
        s.textContent = "Hello, world!"
        document.body.appendChild(s)
    }
})()