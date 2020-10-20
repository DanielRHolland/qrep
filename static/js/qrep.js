var QrepController = {
    getCheckBox: function (resolved) { return resolved ? "check_box" : "check_box_outline_blank";},
    setIssues: function (issues) {
        const printIssue = (issue) => {
            return `
            <li>
               <div class="collection-item">
                  <i class="material-icons">whatshot</i> 
                  ${issue.description} 
                  <a style="cursor:pointer;" class="secondary-content" onclick="QrepController.toggleIssueResolved('${issue.id}')">
                     <i id="${issue.id}resolved" class="material-icons">${QrepController.getCheckBox(issue.resolved)}</i>
                  </a> 
               </div>
            </li>
          `;
        }
        var issuelist = "";
        for (i of issues) {
          issuelist += printIssue(i);
        }
        
        document.getElementById("issues_view").innerHTML =
        `<ul class="collection">
              ${issuelist}
         </ul>
        `;
    },

    displayIssues: function (id) {
        issues = itemissues.find(item => item.id==id).issues
        QrepController.setIssues(issues);
    },

    toggleShow: function (id) {
      var x = document.getElementById(id);
      if (x.style.display === "none") {
        x.style.display = "block";
      } else {
        x.style.display = "none";
      }
    },

    toggleIssueResolved: function (x) {
        console.log(x);
        QrepService.toggleIssueResolved(x);
    },

    setResolved: function (id, resolved) {
        document.getElementById(id+"resolved").innerHTML = QrepController.getCheckBox(resolved)
        M.toast({html: resolved ? "Issue resolved!" : "Issue unresolved."});
    }
}
