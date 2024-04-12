## Files explanation

`article.html` is the exact HTML response from
GET https://medium.com/@andreiboar/fundamentals-of-i-o-in-go-part-2-e7bb68cd5608. Leave it as it is for further
references

`article_formatted.html`is the formatted HTML from `article.html`

`article.json` is the exact response from
GET https://medium.com/andreiboar/fundamentals-of-i-o-in-go-part-2-e7bb68cd5608?format=json. Leave the contents of this
file as it is since is used as a response stub in tests.

`article_formatted.json` contains the formatted JSON without the `])}while(1);</x>` string. Useful to see what
information we have
in the JSON response.