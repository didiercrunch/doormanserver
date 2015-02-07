var bower = require('gulp-bower');
var gulp = require('gulp');
var coffee = require('gulp-coffee');
var help = require('gulp-task-listing');
var clean = require('gulp-clean');
var sass = require('gulp-sass');
var rjs = require("gulp-rjs");


gulp.task('help', help);

gulp.task('bower', function() {
    return bower()
        .pipe(gulp.dest('bower_components/'))
});

gulp.task('sass', ['bower'], function () {
    return gulp.src('./scss/*.scss')
        .pipe(sass({includePaths: ['bower_components/foundation/scss']}))
        .pipe(gulp.dest('./css'));
});

gulp.task('coffee', ['bower'], function() {
    return gulp.src('coffee/**/*.coffee')
        .pipe(coffee({bare: true}))
        .pipe(gulp.dest('js/'))
});

gulp.task('clean', function(){
    return gulp.src(["js/", "bower_components/", "css/"], {read: false})
        .pipe(clean())
});


gulp.task('rjs', ['coffee'], function(){
    return gulp.src('coffee/**/*.coffee')
        .pipe(coffee())
        .pipe(gulp.dest('./dist'))
        .pipe(rjs({baseUrl:'.'}))
});


gulp.task('default', ['bower', 'sass', 'coffee']);



gulp.task('dev', ['bower', 'sass', 'coffee'], function(){
    gulp.watch('coffee/**/*.coffee', ['coffee']);
    gulp.watch('./scss/*.scss', ['sass']);
});
