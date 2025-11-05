# The Playbook
Backend for "[The Playbook](https://github.com/endriym/PickAppBook)" smartphone application. It started as a simple pickupline collector but evolved into a "simple" social.
Made with some sort of basic scalability in mind, so it could theoretically handle a reasonable amount of users for example by sharding elasticsearch.
## About
### Golang Side
My first Golang project, so not everything is optimal. The golang side is composed of:
- **Gin Gonic** as a "little" rest framework
- **Gorm**: as ORM
- **Elasticsearch sdk**: to make my life easier on Elasticsearch
- **Go Mock**: for unit testing using mocks
- **Go Test**: for testing

### Rest of the party
For basic persistence it expects a simple mysql database (for example mariadb) and for all the searching (basically 90% of the project) it uses Elasticsearch 
### Improvements
There's lot to do:
- [ ] All the tests should be revised, most of them are not 100% clear even to me, and some are clearly biased.
- [ ] Add more tests (especially to increase branch coverage)
- [ ] Lots of the elasticsearch queries should be revised, for example the various sortings
- [ ] A more accurate seeder could be written
- [ ] Create utilities scripts for basic document reindexing

And more...
