const parent = document.getElementById('pages');
if (parent && parent.children) {
    highlightCurrentPage(parent);
}

function highlightCurrentPage(parent) {
    Array.from(parent.children).forEach(c => {
        Array.from(c.children).forEach(a => {
            if (a.attributes.href.value === window.location.pathname) {
                a.classList.add('text-primary-600');
            }
        });
    });
}
