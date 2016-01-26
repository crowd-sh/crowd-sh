package us.workmachine;

import android.app.Activity;
import android.os.Bundle;
import com.parse.ui.*;


/**
 * Created by ayerra on 7/1/14.
 */
public class WorkMachineLogin extends Activity {
    @Override
    public void onCreate(Bundle savedInstanceState) {
        ParseLoginBuilder builder = new ParseLoginBuilder(WorkMachine.this);
        startActivityForResult(builder.build(), 0);
    }
}
