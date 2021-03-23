var disqus_config = function () {
    this.page.url = '{{SiteURL}}/{{.posturl}}';
    this.page.identifier = '{{.id}}';
};
(function () {
    var d = document, s = d.createElement('script');
    s.src = 'https://{{DisqusID}}.disqus.com/embed.js';
    s.setAttribute('data-timestamp', +new Date());
    (d.head || d.body).appendChild(s);
})();