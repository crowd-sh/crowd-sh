package us.workmachine;

import android.app.ActionBar;
import android.app.Activity;
import android.content.Intent;
import android.os.Bundle;
import android.util.Log;
import android.view.Menu;
import android.view.MenuItem;
import android.webkit.WebView;

import android.widget.ShareActionProvider;
import com.google.analytics.tracking.android.EasyTracker;

import com.parse.*;

public class WorkMachine extends Activity {

    private ShareActionProvider mShareActionProvider;

    private WebView webView;

    private ParseUser currentUser;

    /**
     * Called when the activity is first created.
     */
    @Override
    public void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.main);

        Parse.initialize(this, "spMy2uicUQgvySv4Pb4sPfYVZE3IfpKQOBTrjc6G", "Nvw4JVfsgpMfk7Jay4u2mrO2OBx4bdvYCvOVRdBP");

        webView = (WebView) findViewById(R.id.webView1);
        webView.clearCache(true);
        webView.clearHistory();
        webView.getSettings().setJavaScriptEnabled(true);
        webView.getSettings().setLoadWithOverviewMode(true);
        webView.getSettings().setUseWideViewPort(true);
        webView.getSettings().setBuiltInZoomControls(true);
        webView.getSettings().setAllowFileAccess(true);
        webView.addJavascriptInterface(new WebAppInterface(this), "Android");
        webView.loadUrl("http://workmachine.us/#/?app=1");

        ParseAnalytics.trackAppOpened(getIntent());

        currentUser = ParseUser.getCurrentUser();
        if (currentUser != null) {

        } else {
            ParseAnonymousUtils.logIn(new LogInCallback() {
                @Override
                public void done(ParseUser user, ParseException e) {
                    if (e != null) {
                        Log.d("WorkMachine", "Anonymous login failed.");
                    } else {
                        Log.d("WorkMachine", "Anonymous user logged in.");
                        currentUser = user;
                    }
                }
            });
        }
    }

    @Override
    public void onStart() {
        super.onStart();
        EasyTracker.getInstance(this).activityStart(this);
    }

    @Override
    public void onStop() {
        super.onStop();
        EasyTracker.getInstance(this).activityStop(this);
    }

    @Override
    public boolean onCreateOptionsMenu(Menu menu) {
        // Inflate menu resource file.
        getMenuInflater().inflate(R.menu.main_menu, menu);

        MenuItem login = menu.findItem(R.id.menu_item_login);
        if(ParseUser.getCurrentUser() != null) {
            login.setEnabled(false);
        }

        // Locate MenuItem with ShareActionProvider
        MenuItem item = menu.findItem(R.id.menu_item_share);

        // Fetch and store ShareActionProvider
        mShareActionProvider = (ShareActionProvider) item.getActionProvider();
        setShareIntent();

        MenuItem feedbackItem = menu.findItem(R.id.menu_item_feedback);
        feedbackItem.setOnMenuItemClickListener(new MenuItem.OnMenuItemClickListener() {
            @Override
            public boolean onMenuItemClick(MenuItem menuItem) {
                final Intent emailIntent = new Intent(android.content.Intent.ACTION_SEND);

                /* Fill it with Data */
                emailIntent.setType("plain/text");
                emailIntent.putExtra(android.content.Intent.EXTRA_EMAIL, new String[]{"abhi@workmachine.us"});
                emailIntent.putExtra(android.content.Intent.EXTRA_SUBJECT, "WorkMachine Feedback");
                emailIntent.putExtra(android.content.Intent.EXTRA_TEXT, "Enter your feedback.");

                /* Send it off to the Activity-Chooser */
                startActivity(Intent.createChooser(emailIntent, "Send mail..."));

                return true;
            }
        });

        // Return true to display menu
        return true;
    }

    // Call to update the share intent
    private void setShareIntent() {
        Intent shareIntent = new Intent(Intent.ACTION_SEND);
        shareIntent.setType("text/plain");
        shareIntent.putExtra(Intent.EXTRA_SUBJECT, "Share WorkMachine");
        shareIntent.putExtra(Intent.EXTRA_TEXT, "Help the world with WorkMachine https://play.google.com/store/apps/details?id=us.workmachine");

        mShareActionProvider.setShareIntent(shareIntent);
    }
}
