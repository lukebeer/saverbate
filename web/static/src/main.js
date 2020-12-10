import 'bootstrap/dist/css/bootstrap.min.css';
import './style.css';

import $jQuery from 'jquery';
import 'bootstrap/js/dist/modal';
import Cookies from 'js-cookie';

if (!cfg.currentUser.loggedIn) {
  $jQuery(function() {
    const overEighteenCookie = Cookies.get('over18');

    if (!overEighteenCookie) {
      $jQuery('.age-confirm-dialog').on('shown.bs.modal', function(){
        $jQuery('.age-confirm-btn').on('click', function(){
          Cookies.set('over18', true);
          $jQuery('.age-confirm-dialog').modal('hide');
        });
      });

      $jQuery('.age-confirm-dialog').modal(
        {
          keyboard: false,
          backdrop: 'static'
        }
      );
    }
  });
}



