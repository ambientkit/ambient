var disqus_config = function () {
    this.page.url = '{{.SiteURL}}{{disqus_PageURL}}';
    //this.page.identifier = '{{.SiteURL}}{{disqus_PageURL}}';
};
(function () {
    var d = document, s = d.createElement('script');
    s.src = 'https://{{.DisqusID}}.disqus.com/embed.js';
    s.setAttribute('data-timestamp', +new Date());
    (d.head || d.body).appendChild(s);
})();