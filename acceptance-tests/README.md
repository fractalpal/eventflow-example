## Acceptance tests
Are designed to test different use cases on a living system.
Could be used inside cluster for automating e2e behaviour in system.

For demo purpose only:
- payment creation
- payment deletion

### Behave
Implementation is done using BDD framework in python. 
(https://behave.readthedocs.io)

```
# instal behave
$ pip install behave
```

### How to run
1. Payment Command and Query are running
2. Run behave:

```
# run tests
$ behave
```
