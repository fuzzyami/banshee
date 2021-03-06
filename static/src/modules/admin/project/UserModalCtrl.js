/*@ngInject*/
module.exports = function ($scope, $mdDialog, $stateParams, $translate, toastr, Project, params) {
  $scope.titles = {
    addUserToProject: $translate.instant('ADMIN_USER_ADD_TEXT')
  };
  $scope.opt = params.opt;

  $scope.users = params.users;

  $scope.autoComplete = {
    searchText: ''
  };


  $scope.cancel = function() {
    $mdDialog.cancel();
  };

  $scope.submit = function() {
    Project.addUserToProject({
      id: $stateParams.id,
      name: $scope.autoComplete.selectedItem.name
    }).$promise
    .then(function() {
      $mdDialog.hide($scope.autoComplete.selectedItem);
    })
    .catch(function(err) {
      toastr.error(err.msg);
    });
  };
};
