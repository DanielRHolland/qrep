var QrepController = {
    getCheckBox: function (resolved) { return resolved ? "check_box" : "check_box_outline_blank";},
    setIssues: function (item) {
        issues = item.issues
        const printIssue = (issue) => {
            return `
            <li>
               <div class="collection-item">
                  <i class="material-icons">whatshot</i> 
                  ${issue.description} 
                  <a style="cursor:pointer;" class="secondary-content" onclick="QrepController.toggleIssueResolved('${issue.id}')">
                     <i class="material-icons ${issue.id}resolved" >${QrepController.getCheckBox(issue.resolved)}</i>
                  </a> 
               </div>
            </li>
          `;
        }
        var issuelist = "";
        for (i of issues) {
          issuelist += printIssue(i);
        }
        var htmlToSet = `
        <h4> ${item.name} </h4>
        <ul class="collection">
              ${issuelist}
         </ul>
        `;
        document.getElementById("issues_modal_view").innerHTML = htmlToSet;
        document.getElementById("issues_view").innerHTML = htmlToSet;

        let isMobile = window.matchMedia("only screen and (max-width: 760px)").matches;
        if (isMobile){
            let elem = document.getElementById("issuesModal");
            M.Modal.getInstance(elem).open();
        }     
    },

    displayIssues: function (id) {
        item = itemissues.find(item => item.id==id);
        QrepController.setIssues(item);
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
