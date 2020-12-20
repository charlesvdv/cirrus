import unittest

import testutils

class IntegrationTestCase(unittest.TestCase):
    def __init__(self, name: str, client: testutils.Client):
        super().__init__(name)
        self._client = client


class IntegrationTestLoader(unittest.TestLoader):
    def __init__(self, client: testutils.Client, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self._client = client

    # See https://stackoverflow.com/a/37916673 for more info.
    # The basic idea is to inject dependencies to `IntegrationTestCase` when creating
    # the test case.
    def loadTestsFromTestCase(self, testCaseClass, **kwargs):
        if not issubclass(testCaseClass, IntegrationTestCase):
            raise 'Test case should be derived from IntegrationTestCase'

        testCaseNames = self.getTestCaseNames(testCaseClass)
        if not testCaseNames and hasattr(testCaseClass, 'runTest'):
            testCaseNames = ['runTest']

        finalTestCases = []
        for testCaseName in testCaseNames:
            finalTestCases.append(testCaseClass(testCaseName, self._client))
        
        return self.suiteClass(finalTestCases)