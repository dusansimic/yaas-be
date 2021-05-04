const {task, src, dest, series} = require('gulp');
const minify = require('gulp-minify');

task('js', () => {
	return src('yaas.js')
		.pipe(minify({noSource: true}))
		.pipe(dest('.'));
});

task('default', series('js'));
