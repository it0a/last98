/*global window, document, jQuery, browser, $*/
"use strict";
$(document).on('change', '.btn-file :file', function() {
    var input = $(this),
    numFiles = input.get(0).files ? input.get(0).files.length : 1,
    label = input.val().replace(/\\/g, '/').replace(/.*\//, '');
        input.trigger('fileselect', [numFiles, label]);
});
$(document).ready( function() {
    $('.btn-file :file').on('fileselect', function(event, numFiles, label) {
        /*jslint unparam:true*/
        $('input[type="file"]').closest("div").find('input[type="text"]').prop("value", label);
    });
});
