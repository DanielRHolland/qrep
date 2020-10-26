var QrepService = {
    test: (x) => {console.log(x)},

    putIssue: async (issueId) => {  
      var i = QrepService.findIssue(itemissues, issueId);
      const response = await fetch("/issue/"+issueId, {
        method: 'PUT',
        body: JSON.stringify(i), // string or object
        headers: {
          'Content-Type': 'application/json'
        }
      });
      const myJson = await response.json(); //extract JSON from the http response
      // do something with myJson
      console.log(myJson);
    },

    toggleIssueResolved: async (issueId) => {
      var i = QrepService.findIssue(itemissues, issueId);
      i.resolved = !i.resolved;
      const response = await fetch("/issue/" + issueId , {
        method: 'PUT',
        body: JSON.stringify(i), // string or object
        headers: {
          'Content-Type': 'application/json'
        }
     });
      const myJson = await response.json();
      console.log(myJson);
      QrepController.setResolved(i.id, i.resolved)
    },

    findIssue: (itemissues, issueId) => {
      let issueIdMatch = (issue) => issue.id==issueId;
      return itemissues
        .find(item => item.issues.some(issueIdMatch))
        .issues.find(issueIdMatch);
    },

    deleteSelectedItems: () => {
      let form = document.getElementById("itemsselector");
      form.action="remove";
      form.submit();
    }
}
