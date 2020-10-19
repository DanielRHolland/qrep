
function setIssues(issues) {
    const printIssue = (issue) => {
        return `<li>
           <div class="collection-item"><i class="material-icons">whatshot</i> ${issue.description} </div>
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
}


function displayIssues(id) {
    issues = itemissues.find(item => item.id==id).issues
    setIssues(issues);
}


function toggleShow(id) {
  var x = document.getElementById(id);
  if (x.style.display === "none") {
    x.style.display = "block";
  } else {
    x.style.display = "none";
  }
}
