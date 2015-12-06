function TaskCtrl($scope, $http) {
    $scope.tasks = [];
    $scope.working = false;

    var logError = function(data, status) {
        console.log('code '+status+': '+data);
        $scope.working = false;
    };

    var refresh = function() {
        return $http.get('/api/todos').
            success(function(data) { $scope.tasks = data; console.log(data) }).
            error(logError);
    };

    $scope.addTodo = function() {
        $scope.working = true;
        $http.post('/api/todos', {Title: $scope.todoText}).
            error(logError).
            success(function() {
                refresh().then(function() {
                    $scope.working = false;
                    $scope.todoText = '';
                });
            });
    };

    $scope.toggleDone = function(task) {
        data = {Id: task.Id, Title: task.Title, Done: !task.Done};
        $http.put('/api/todos/'+task.Id, data).
            error(logError).
            success(function() { task.Done = !task.Done; });
    };

    refresh().then(function() { $scope.working = false; });
}
