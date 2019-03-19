### Integration tests
Are designed to test integration of payment (command) service with database.
It spins new server and makes some request during tests. 

For demo purpose only:
- payment creation
- payment deletion

#### BDD testing
Integration test are created using BDD testing frameworks:
- [ginkgo](https://github.com/onsi/ginkgo)
- [gomega](https://github.com/onsi/gomega)
 