var QrepController = {
    getCheckBox: function (resolved) { return resolved ? "check_box" : "check_box_outline_blank";},
    displayIssues: function (id, showResolvedIssues=false) {
        let item = itemissues.find(item => item.id==id);
        issues = item.issues
        const printIssue = (issue) => {
            return `
            <li>
               <div class="collection-item">
                  <span>${issue.description}</span> 
                  <a style="cursor:pointer;" class="secondary-content" onclick="QrepController.toggleIssueResolved('${issue.id}')">
                     <i id="${issue.id}resolved" class="material-icons">${QrepController.getCheckBox(issue.resolved)}</i>
                  </a> 
               </div>
            </li>
          `;
        }
        var issuelist = "";
        for (i of issues) {
          if (i.resolved == showResolvedIssues) issuelist += printIssue(i);
        }
//        var toggleShowResolvedButtonText = 
        var toggleViewResolvedButton =`
            <a class="btn-flat waves-effect waves-light" onclick="QrepController.displayIssues('${id}',${!showResolvedIssues})">
                <i class="material-icons">compare_arrows</i>
                ${showResolvedIssues ? "SHOW UNRESOLVED ISSUES"   :  "SHOW RESOLVED ISSUES"}
            </a>`

        var htmlToSet =
        `<h4> ${item.name} </h4>
         ${toggleViewResolvedButton}
         <ul class="collection">
              ${issuelist}
         </ul>
        `;
        document.getElementById("issues_view").innerHTML = htmlToSet;
        document.getElementById("issues_modal_view").innerHTML = htmlToSet;
        let isMobile = window.matchMedia("only screen and (max-width: 760px)").matches;
        if (isMobile){
            let elem = document.getElementById("issuesModal");
            M.Modal.getInstance(elem).open();
        } 
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
        let elems = document.getElementsByClassName(id+"resolved");
        for (i=0; i<elems.length; i++) {
            elems[i].innerHTML = QrepController.getCheckBox(resolved);
        }
        M.toast({html: resolved ? "Issue resolved!" : "Issue unresolved."});
    }
}
