"use strict";

import gulp from "gulp";
import sass from "gulp-sass";
import uglify from "gulp-uglify";
import browserify from "browserify";
import source from "vinyl-source-stream";

gulp.task("css", () => {
  return gulp
    .src("css/tracklog.scss")
    .pipe(sass.sync())
    .pipe(gulp.dest("public"));
});

gulp.task("js", () => {
  return browserify({ entries: ["js/tracklog.js"], debug: true })
    .transform("babelify")
    .bundle()
    .pipe(source("tracklog.js"))
    .pipe(gulp.dest("public"));
});

gulp.task("compress", () => {
  return gulp
    .src("public/tracklog.js")
    .pipe(uglify())
    .pipe(gulp.dest("public"));
});

gulp.task("watch", () => {
  gulp.watch("js/**/*.js", ["js"]);
  gulp.watch("css/**/*.scss", ["css"]);
});

gulp.task("default", ["css", "js"]);
