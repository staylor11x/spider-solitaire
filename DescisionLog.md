# Log of Key Decisions in the Project

### Testing Predicimant [18/09/2025]

We are facing issues with testing internal/external methods, most of these issues seem to be stemming from the test tools library that we implemented. When we implemented this it caused us to have circular dependencies in our tests, we fixed this by moving the tests into a "different" package. 

The problem that this caused is that now we can only test "exposed" methods that are part of the public facing API, originally this was fine but as we started to introduce more complicated logic into the un-exported helper functions this issue quickly surfaced.

I am making the decision based on this to scrap the external test tools module and either simply include the helper functions for the tests within the test file they are required in or not use them at all.

I think this idea of test tool helper functions works better in Java esq languages 