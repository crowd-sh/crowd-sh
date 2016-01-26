package us.workmachine;

import android.content.Context;
import android.widget.Toast;
import com.parse.ParseObject;
import com.parse.ParseUser;

/**
 * Created by ayerra on 1/6/14.
 */
public class WebAppInterface {
    Context mContext;

    /** Instantiate the interface and set the context */
    WebAppInterface(Context c) {
        mContext = c;
    }

    /** Show a toast from the web page */
    // @JavascriptInterface
    public void incrementTasks() {
        if(ParseUser.getCurrentUser() != null) {
            ParseUser.getCurrentUser().increment("WorkCount");
            ParseUser.getCurrentUser().saveInBackground();
        }
    }

    // @JavascriptInterface
    public String getLeaderboard() {

        return "123";
    }
}
